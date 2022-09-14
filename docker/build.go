package docker

import (
	"context"
	"fmt"
	"github.com/brunoscheufler/atlas/atlasfile"
	"github.com/brunoscheufler/atlas/exec"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
)

func BuildArtifact(ctx context.Context, logger logrus.FieldLogger, artifact *atlasfile.ArtifactConfig) error {
	logger.WithField("artifact", artifact.Name).Infoln("Building artifact")

	dir := filepath.Dir(artifact.GetDirpath())
	if artifact.Build.Context != "" {
		dir = filepath.Join(dir, artifact.Build.Context)
	}

	logger.WithField("dir", dir).Infoln("Building artifact")

	imageName := atlasfile.BuildImageName(artifact)

	args := []string{
		"build",
		"-t",
		imageName,
	}

	if artifact.Build.Dockerfile != "" {
		dockerfilePath := filepath.Join(dir, artifact.Build.Dockerfile)
		args = append(args, "-f", dockerfilePath)
	}

	if artifact.Build.BuildArgs != nil {
		for key, value := range artifact.Build.BuildArgs {
			args = append(args, "--build-arg", fmt.Sprintf("%s=%s", key, value))
		}
	}

	if artifact.Build.Target != "" {
		args = append(args, "--target", artifact.Build.Target)
	}

	args = append(args, dir)

	err := exec.RunCommand(ctx, logger, fmt.Sprintf("docker %s", strings.Join(args, " ")), dir, nil)
	if err != nil {
		return fmt.Errorf("could not build artifact %s: %w", artifact.Name, err)
	}

	return nil
}
