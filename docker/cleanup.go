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
