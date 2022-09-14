package docker

import (
	"context"
	"fmt"
	"github.com/brunoscheufler/atlas/exec"
	"github.com/sirupsen/logrus"
)

func CreateVolume(ctx context.Context, logger logrus.FieldLogger, name string) error {
	err := exec.RunCommand(ctx, logger, fmt.Sprintf("docker volume create %s", name), "", nil)
	if err != nil {
		return fmt.Errorf("could not create volume %s: %w", name, err)
	}

	return nil
}
