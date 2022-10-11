package atlas

import (
	"context"
	"fmt"
	"github.com/brunoscheufler/atlas/atlasfile"
	"github.com/brunoscheufler/atlas/docker"
	"github.com/sirupsen/logrus"
)

func Stop(ctx context.Context, logger logrus.FieldLogger, version, cwd string, stackName, serviceName string) error {
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

	stack := statefile.GetStack(stackName)
	if stack == nil {
		return fmt.Errorf("stack %s not found", stackName)
	}

	service := stack.GetService(serviceName)
	if service == nil {
		return fmt.Errorf("service %s not found", serviceName)
	}

	if service.ContainerInfos.State == "exited" {
		logger.Infof("Service %s is already stopped", serviceName)
		return nil
	}

	err = docker.StopContainer(ctx, service.ContainerName)
	if err != nil {
		return fmt.Errorf("could not stop container: %w", err)
	}

	return nil
}
