package atlas

import (
	"context"
	"fmt"
	"github.com/brunoscheufler/atlas/atlasfile"
	"github.com/brunoscheufler/atlas/docker"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
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

	stateFile, err := readState(ctx, cwd, version, logger)
	if err != nil {
		return fmt.Errorf("could not read state file: %w", err)
	}

	if stateFile == nil {
		logger.Infoln("No state file found, nothing to do")
		return nil
	}

	stateFileStacks, err := stateFile.GetStacks(stackNames)
	if err != nil {
		return fmt.Errorf("could not get stacks: %w", err)
	}

	for _, stack := range stateFileStacks {
		logger.WithField("stack", stack.Name).WithField("network", stack.Network).Infof("- Stopping stack %s\n", stack.Name)

		{
			g, ctx := errgroup.WithContext(ctx)
			for _, service := range stack.Services {
				service := service
				g.Go(func() error {
					logger.WithFields(logrus.Fields{
						"stack":   stack.Name,
						"service": service.Name,
					}).Infof("\t- Stopping service %s\n", service.Name)

					err = docker.DeleteContainer(ctx, logger, service.ContainerName)
					if err != nil {
						return fmt.Errorf("could not stop service: %w", err)
					}

					return nil
				})
			}

			err = g.Wait()
			if err != nil {
				return fmt.Errorf("could not stop stack services: %w", err)
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
