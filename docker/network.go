package docker

import (
	"context"
	"fmt"
	"github.com/brunoscheufler/atlas/exec"
	"github.com/sirupsen/logrus"
)

func CreateNetwork(ctx context.Context, logger logrus.FieldLogger, name string) error {
	err := exec.RunCommand(ctx, logger, fmt.Sprintf("docker network create %s", name), "", nil, false)
	if err != nil {
		return fmt.Errorf("could not create network %s: %w", name, err)
	}

	return nil
}
