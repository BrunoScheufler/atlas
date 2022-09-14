package atlasfile

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/brunoscheufler/atlas/exec"
	"github.com/brunoscheufler/atlas/helper"
	"github.com/brunoscheufler/atlas/protobuf"
	"github.com/cenkalti/backoff/v4"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

const JumpUpLimit = 5

func FindRootDir(cwd string) (string, error) {
	jumpedUp := 0

	for {
		if jumpedUp > JumpUpLimit {
			return "", fmt.Errorf("could not find root directory")
		}

		// Check in cwd if .atlas exists as Atlasfile.root.*
		atlasDirPath := filepath.Join(cwd, ".atlas")

		if helper.FileExists(filepath.Join(atlasDirPath, "Atlasfile.root.go")) || helper.FileExists(filepath.Join(atlasDirPath, "Atlasfile.root.ts")) {
			return cwd, nil
		}

		// Go up one directory
		cwd = filepath.Dir(cwd)
		jumpedUp++
	}
}

func FindAtlasDirectories(dir string) ([]string, error) {
	files := make([]string, 0)

	skipDirs := []string{
		"node_modules",
	}

	shouldBeSkipped := func(path string) bool {
		for _, skipDir := range skipDirs {
			if filepath.Base(path) == skipDir {
				return true
			}
		}
		return false
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && shouldBeSkipped(path) {
			return filepath.SkipDir
		}

		if info.IsDir() && info.Name() == ".atlas" {
			files = append(files, path)
			return filepath.SkipDir
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not walk directory: %w", err)
	}

	return files, nil
}

func Collect(ctx context.Context, logger logrus.FieldLogger, cwd string) (*Atlasfile, error) {
	cwd, err := FindRootDir(cwd)
	if err != nil {
		return nil, fmt.Errorf("could not find root directory: %w", err)
	}

	// Find all .atlas directories with glob
	paths, err := FindAtlasDirectories(cwd)
	if err != nil {
		return nil, fmt.Errorf("could not glob for .atlas directories: %w", err)
	}

	logger.WithField("paths", paths).Debugln("globbed .atlas directories")

	collectedFiles := make([]Atlasfile, len(paths))

	bar := progressbar.NewOptions(len(paths), progressbar.OptionSetDescription("Reading Atlasfiles"), progressbar.OptionClearOnFinish())

	g, ctx := errgroup.WithContext(ctx)

	for i, path := range paths {
		relpath, err := filepath.Rel(cwd, path)
		if err != nil {
			return nil, fmt.Errorf("could not find relative path: %w", err)
		}

		i, path, bar, relpath := i, path, bar, relpath // https://golang.org/doc/faq#closures_and_goroutines

		g.Go(func() error {
			bar.Describe(fmt.Sprintf("Reading Atlasfile %s", relpath))

			file, err := readAtlasFile(ctx, logger, path)
			if err != nil {
				return fmt.Errorf("could not read .atlas file: %w", err)
			}
			collectedFiles[i] = *file
			_ = bar.Add(1)
			return nil
		})

	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	_ = bar.Clear()
	_ = bar.Close()

	return mergeAtlasFiles(collectedFiles), nil
}

// TODO Cache files in .atlas directory and return previous value if unchanged
// TODO Support non-code/Toml Atlasfile
func readAtlasFile(ctx context.Context, logger logrus.FieldLogger, atlasDirPath string) (*Atlasfile, error) {
	// Check if go.mod exists
	if helper.FileExists(filepath.Join(atlasDirPath, "go.mod")) {
		return readGoAtlasFile(ctx, logger, atlasDirPath)
	}

	// Check if package.json exists
	if helper.FileExists(filepath.Join(atlasDirPath, "package.json")) {
		return readTypeScriptAtlasFile(ctx, atlasDirPath)
	}

	return nil, fmt.Errorf("missing go.mod or package.json, cannot infer language to use")
}

func readGoAtlasFile(ctx context.Context, logger logrus.FieldLogger, atlasDirPath string) (*Atlasfile, error) {
	// Check if building the file works
	err := exec.RunCommand(ctx, logger, "go build -o /dev/null .", atlasDirPath, nil)
	if err != nil {
		return nil, fmt.Errorf("could not build atlas file (%s): %w", atlasDirPath, err)
	}

	port, err := helper.FreePort()
	if err != nil {
		return nil, fmt.Errorf("could not find free port: %w", err)
	}

	// Start process in background with go run . in the directory
	cmd, err := exec.StartCommand(ctx, logger, "go run .", atlasDirPath, []string{fmt.Sprintf("PORT=%d", port)})

	backOff := &backoff.ExponentialBackOff{
		InitialInterval:     time.Millisecond * 10,
		MaxInterval:         time.Second * 3,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          5,
		MaxElapsedTime:      0,
		Stop:                backoff.Stop,
		Clock:               backoff.SystemClock,
	}
	backOff.Reset()

	var client protobuf.AtlasfileClient
	var conn *grpc.ClientConn

	// Wait until started up
	attempts := 0
	for {
		select {
		// Either the ctx is canceled
		case <-ctx.Done():
			return nil, fmt.Errorf("could not connect to atlasfile, context canceled")

			// Or the backOff elapses
		case <-time.After(backOff.NextBackOff()):
			// after which we might have to return as we reached the max num of attempts
			if attempts > 10 {
				return nil, fmt.Errorf("could not connect to atlasfile: %w", err)
			}

			// Connect to gRPC endpoint
			conn, err = grpc.DialContext(ctx, fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				attempts++
				continue
			}

			client = protobuf.NewAtlasfileClient(conn)

			// try pinging the endpoint
			_, err = client.Ping(ctx, &protobuf.PingRequest{})
			if err != nil {
				attempts++
				continue
			}
		}

		break
	}

	// Send request to get atlasfile
	res, err := client.Eval(ctx, &protobuf.EvalRequest{})
	if err != nil {
		return nil, fmt.Errorf("could not eval atlasfile: %w", err)
	}

	var atlasfile Atlasfile
	err = json.Unmarshal([]byte(res.Output), &atlasfile)
	if err != nil {
		return nil, fmt.Errorf("could not parse atlasfile: %w", err)
	}
	atlasfile.dirpath = atlasDirPath

	_ = conn.Close()

	// Shut down process
	err = cmd.Process.Signal(syscall.SIGTERM)
	if err != nil {
		return nil, fmt.Errorf("could not shut down atlasfile provider: %w", err)
	}

	go func() {
		select {
		case <-ctx.Done():
		case <-time.After(time.Second * 5):
			_ = cmd.Process.Signal(syscall.SIGKILL)
		}
	}()

	return &atlasfile, nil
}

func readTypeScriptAtlasFile(ctx context.Context, atlasDirPath string) (*Atlasfile, error) {
	return nil, fmt.Errorf("unsupported")
}

func mergeAtlasFiles(files []Atlasfile) *Atlasfile {
	final := &Atlasfile{
		Artifacts: make([]ArtifactConfig, 0),
		Services:  make([]ServiceConfig, 0),
		Stacks:    make([]StackConfig, 0),
	}

	for _, file := range files {
		for _, artifact := range file.Artifacts {
			artifact.dirpath = file.dirpath
			final.Artifacts = append(final.Artifacts, artifact)
		}

		for _, service := range file.Services {
			service.dirpath = file.dirpath
			// Move artifact from service scope into artifacts
			if service.Artifact != nil && service.Artifact.Artifact != nil {
				// TODO make sure no artifact with same name exists
				svcArtifact := *service.Artifact.Artifact
				svcArtifact.dirpath = file.dirpath
				final.Artifacts = append(final.Artifacts, svcArtifact)
				service.Artifact = &ArtifactRef{Name: svcArtifact.Name}
			}
			final.Services = append(final.Services, service)
		}

		for _, stack := range file.Stacks {
			stack.dirpath = file.dirpath
			final.Stacks = append(final.Stacks, stack)
		}
	}

	return final
}
