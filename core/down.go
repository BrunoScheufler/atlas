package atlas

import (
	"context"
	"fmt"
	"github.com/brunoscheufler/atlas/atlasfile"
	"github.com/brunoscheufler/atlas/docker"
	"github.com/sirupsen/logrus"
)

func Down(ctx context.Context, logger logrus.FieldLogger, cwd, version string, stackNames []string, all bool) error {
	if !docker.IsRunning(ctx) {
		return fmt.Errorf("docker is not running")
	}

	if all {
		err := docker.CleanupAll(ctx, logger)
		if err != nil {
			return fmt.Errorf("could not cleanup: %w", err)
		}

		return nil
	}

	cwd, err := atlasfile.FindRootDir(cwd)
	if err != nil {
		return fmt.Errorf("could not find root directory: %w", err)
	}

	stateFile, err := readState(cwd, version, logger)
	if err != nil {
		return fmt.Errorf("could not read state file: %w", err)
	}

	if stateFile == nil {
		logger.Infoln("No state file found, nothing to do")
		return nil
	}

	var stacks []StateStack
	if len(stackNames) == 0 {
		stacks = stateFile.Stacks
	} else {
		for _, stackName := range stackNames {
			for i, stack := range stateFile.Stacks {
				if stack.Name == stackName {
					stacks = append(stacks, stateFile.Stacks[i])
					break
				}
			}

			return fmt.Errorf("could not find stack %q in state ", stackName)
		}
	}

	for _, stack := range stateFile.Stacks {
		for _, service := range stack.Services {
			logger.WithFields(logrus.Fields{
				"stack":   stack.Name,
				"service": service.Name,
			}).Debugf("Stopping service")

			err = docker.DeleteContainer(ctx, logger, service.ContainerName)
			if err != nil {
				return fmt.Errorf("could not stop service: %w", err)
			}
		}

		err = docker.DeleteNetwork(ctx, logger, stack.Network)
		if err != nil {
			return fmt.Errorf("could not delete network: %w", err)
		}
	}

	for _, volumeName := range stateFile.Volumes {
		err = docker.DeleteVolume(ctx, logger, volumeName)
		if err != nil {
			return fmt.Errorf("could not delete volume: %w", err)
		}
	}

	err = clearStatefile(cwd)
	if err != nil {
		return fmt.Errorf("could not clear state file: %w", err)
	}

	return nil
}
