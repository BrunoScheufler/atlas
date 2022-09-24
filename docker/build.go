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

func BuildArtifact(ctx context.Context, logger logrus.FieldLogger, artifact *atlasfile.ArtifactConfig, rootDir string) error {
	artifactDir := filepath.Dir(artifact.GetDirpath())
	if artifact.Build.Context != "" {
		artifactDir = filepath.Join(artifactDir, artifact.Build.Context)
	}

	logger.WithFields(logrus.Fields{
		"artifact": artifact.Name,
		"rootDir":  rootDir,
		"dirpath":  artifactDir,
		"context":  artifact.Build.Context,
	}).Debugf("Building artifact")

	logger.WithField("dir", artifactDir).WithField("artifact", artifact.Name).Infoln("Building artifact")

	imageName := atlasfile.BuildImageName(artifact)

	args := []string{
		"build",
		"-t",
		imageName,
	}

	if artifact.Build.Dockerfile != "" {
		dockerfilePath := filepath.Join(artifactDir, artifact.Build.Dockerfile)
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

	args = append(args, artifactDir)

	err := exec.RunCommand(ctx, logger, fmt.Sprintf("docker %s", strings.Join(args, " ")), artifactDir, nil)
	if err != nil {
		return fmt.Errorf("could not build artifact %s: %w", artifact.Name, err)
	}

	return nil
}
