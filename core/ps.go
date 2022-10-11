package atlas

import (
	"context"
	"fmt"
	"github.com/brunoscheufler/atlas/atlasfile"
	"github.com/brunoscheufler/atlas/docker"
	"github.com/sirupsen/logrus"
)

func Ps(ctx context.Context, logger logrus.FieldLogger, cwd, version string, stackNames []string) error {
	cwd, err := atlasfile.FindRootDir(cwd)
	if err != nil {
		return fmt.Errorf("could not find root directory: %w", err)
	}

	if !docker.IsRunning(ctx) {
		return fmt.Errorf("docker is not running")
	}

	statefile, err := readState(ctx, cwd, version, logger)
	if err != nil {
		return fmt.Errorf("could not read state file: %w", err)
	}

	if statefile == nil {
		logger.Infoln("No state file found, nothing to do")
		return nil
	}

	stacks, err := statefile.GetStacks(stackNames)
	if err != nil {
		return fmt.Errorf("could not get stacks: %w", err)
	}

	for _, stack := range stacks {
		logger.WithField("stack", stack.Name).WithField("network", stack.Network).Infof("%s (%d):", stack.Name, len(stack.Services))
		for _, service := range stack.Services {
			logger.WithField("containerName", service.ContainerName).Infof("\t- %s: %s (%s)", service.Name, service.ContainerInfos.Status, service.ContainerInfos.State)
		}
	}

	return nil
}
