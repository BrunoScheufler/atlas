package docker

import (
	"context"
	"github.com/brunoscheufler/atlas/exec"
)

func IsRunning(ctx context.Context) bool {
	err := exec.RunCommand(ctx, nil, "docker info", "", nil)
	if err != nil {
		return false
	}

	return true
}
