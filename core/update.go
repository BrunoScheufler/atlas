package atlas

import (
	"context"
	"fmt"
	"github.com/brunoscheufler/atlas/atlasfile"
	"github.com/brunoscheufler/atlas/exec"
	"github.com/brunoscheufler/atlas/helper"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

func Update(ctx context.Context, logger logrus.FieldLogger, cwd string) error {
	cwd, err := atlasfile.FindRootDir(cwd)
	if err != nil {
		return fmt.Errorf("could not find root directory: %w", err)
	}

	// Find all .atlas directories with glob
	paths, err := atlasfile.FindAtlasDirectories(cwd)
	if err != nil {
		return fmt.Errorf("could not glob for .atlas directories: %w", err)
	}

	bar := progressbar.NewOptions(len(paths), progressbar.OptionSetDescription("Updating Atlasfiles"), progressbar.OptionClearOnFinish())
	defer func() {
		_ = bar.Finish()
		_ = bar.Clear()
		_ = bar.Close()
	}()

	for _, path := range paths {
		relPath, err := filepath.Rel(cwd, path)
		if err != nil {
			return fmt.Errorf("could not get relative path for %q: %w", path, err)
		}

		// Check if go.mod exists
		if helper.FileExists(filepath.Join(path, "go.mod")) {
			bar.Describe(fmt.Sprintf("Updating Go Atlasfile (%s)", relPath))

			// Run go get -u && go mo tidy
			err := exec.RunCommand(ctx, logger, "go get -u", path, nil)
			if err != nil {
				return fmt.Errorf("could not run go get -u: %w", err)
			}

			err = exec.RunCommand(ctx, logger, "go mod tidy", path, nil)
			if err != nil {
				return fmt.Errorf("could not run go mod tidy: %w", err)
			}

			_ = bar.Add(1)

			continue
		}

		// Check if package.json exists
		if helper.FileExists(filepath.Join(path, "package.json")) {
			return fmt.Errorf("not implemented yet")
		}
	}

	logger.Infof("Updated %d Atlasfiles", len(paths))

	return nil
}
