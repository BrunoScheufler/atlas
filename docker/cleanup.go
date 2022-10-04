package docker

import (
	"context"
	"fmt"
	"github.com/brunoscheufler/atlas/exec"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
)

func CleanupAll(ctx context.Context, logger logrus.FieldLogger) error {
	bar := progressbar.NewOptions(3, progressbar.OptionSetDescription("Cleaning up"), progressbar.OptionClearOnFinish())
	defer func() {
		_ = bar.Finish()
		_ = bar.Clear()
		_ = bar.Close()
	}()

	bar.Describe("Cleaning up containers")
	err := exec.RunCommand(ctx, logger, "docker container stop $(docker container ls -q --filter Name=atlas-) || true", exec.RunCommandOptions{})
	if err != nil {
		return fmt.Errorf("could not remove containers: %w", err)
	}

	err = exec.RunCommand(ctx, logger, "docker container rm -f -v $(docker container ls -a -q --filter Name=atlas-) || true", exec.RunCommandOptions{})
	if err != nil {
		return fmt.Errorf("could not remove containers: %w", err)
	}
	_ = bar.Add(1)

	bar.Describe("Cleaning up volumes")
	err = exec.RunCommand(ctx, logger, "docker volume rm -f $(docker volume ls -q --filter Name=atlas-) || true", exec.RunCommandOptions{})
	if err != nil {
		return fmt.Errorf("could not remove volumes: %w", err)
	}
	_ = bar.Add(1)

	bar.Describe("Cleaning up networks")
	err = exec.RunCommand(ctx, logger, "docker network rm $(docker network ls -q --filter Name=atlas-) || true", exec.RunCommandOptions{})
	if err != nil {
		return fmt.Errorf("could not remove networks: %w", err)
	}
	_ = bar.Add(1)

	return nil
}
