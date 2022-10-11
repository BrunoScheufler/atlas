package docker

import (
	"context"
	"fmt"
	"github.com/brunoscheufler/atlas/exec"
	"github.com/sirupsen/logrus"
)

func CleanupAll(ctx context.Context, logger logrus.FieldLogger) error {

	logger.Infoln("Cleaning up containers")
	err := exec.RunCommand(ctx, logger, "docker container stop $(docker container ls -q --filter Name=atlas-) || true", exec.RunCommandOptions{})
	if err != nil {
		return fmt.Errorf("could not remove containers: %w", err)
	}

	err = exec.RunCommand(ctx, logger, "docker container rm -f -v $(docker container ls -a -q --filter Name=atlas-) || true", exec.RunCommandOptions{})
	if err != nil {
		return fmt.Errorf("could not remove containers: %w", err)
	}

	logger.Infoln("Cleaning up volumes")
	err = exec.RunCommand(ctx, logger, "docker volume rm -f $(docker volume ls -q --filter Name=atlas-) || true", exec.RunCommandOptions{})
	if err != nil {
		return fmt.Errorf("could not remove volumes: %w", err)
	}

	logger.Infoln("Cleaning up networks")
	err = exec.RunCommand(ctx, logger, "docker network rm $(docker network ls -q --filter Name=atlas-) || true", exec.RunCommandOptions{})
	if err != nil {
		return fmt.Errorf("could not remove networks: %w", err)
	}

	return nil
}

func DeleteContainer(ctx context.Context, logger logrus.FieldLogger, containerName string) error {
	existingContainer, err := GetContainerInfo(ctx, containerName)
	if err != nil {
		return fmt.Errorf("could not get container info: %w", err)
	}

	if existingContainer == nil {
		return nil
	}

	err = exec.RunCommand(ctx, logger, fmt.Sprintf("docker container stop %s && docker container rm -f -v %s", containerName, containerName), exec.RunCommandOptions{})
	if err != nil {
		return fmt.Errorf("could not delete network: %w", err)
	}

	return nil
}

func DeleteNetwork(ctx context.Context, logger logrus.FieldLogger, networkName string) error {
	existingId, err := GetNetworkId(ctx, networkName)
	if err != nil {
		return fmt.Errorf("could not get network id: %w", err)
	}

	if existingId == "" {
		return nil
	}

	err = exec.RunCommand(ctx, logger, fmt.Sprintf("docker network rm %s", networkName), exec.RunCommandOptions{})
	if err != nil {
		return fmt.Errorf("could not delete network: %w", err)
	}

	return nil
}

func DeleteVolume(ctx context.Context, logger logrus.FieldLogger, volumeName string) error {
	err := exec.RunCommand(ctx, logger, fmt.Sprintf("docker volume rm -f %s", volumeName), exec.RunCommandOptions{})
	if err != nil {
		return fmt.Errorf("could not delete volume %s: %w", volumeName, err)
	}

	return nil
}
