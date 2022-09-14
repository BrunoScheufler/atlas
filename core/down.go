package atlas

import (
	"context"
	"fmt"
	"github.com/brunoscheufler/atlas/docker"
	"github.com/sirupsen/logrus"
)

func Down(ctx context.Context, logger logrus.FieldLogger, cwd string, stackNames []string) error {
	// TODO only drop workloads from current stacks (and maybe orphans, too?)
	err := docker.CleanupAll(ctx, logger)
	if err != nil {
		return fmt.Errorf("could not cleanup: %w", err)
	}

	return nil
}
