package docker

import (
	"context"
	"github.com/brunoscheufler/atlas/exec"
)

func IsRunning(ctx context.Context) bool {
	err := exec.RunCommand(ctx, nil, "docker info", exec.RunCommandOptions{})
	if err != nil {
		return false
	}

	return true
}
