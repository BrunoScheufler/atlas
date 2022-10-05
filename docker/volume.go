package docker

import (
	"context"
	"fmt"
	"github.com/brunoscheufler/atlas/atlasfile"
	"github.com/brunoscheufler/atlas/exec"
	"github.com/brunoscheufler/atlas/helper"
	"github.com/sirupsen/logrus"
)

func CreateVolume(ctx context.Context, logger logrus.FieldLogger, name string) error {
	err := exec.RunCommand(ctx, logger, fmt.Sprintf("docker volume create %s", name), exec.RunCommandOptions{})
	if err != nil {
		return fmt.Errorf("could not create volume %s: %w", name, err)
	}

	return nil
}

type EnsuredVolume struct {
	Stack        string
	Service      string
	VolumeName   string
	PhysicalName string
}

type EnsuredVolumes []EnsuredVolume

func (e *EnsuredVolumes) Get(stackName, serviceName, volumeName string) string {
	for _, vol := range *e {
		if vol.Stack == stackName && vol.Service == serviceName && vol.VolumeName == volumeName {
			return vol.PhysicalName
		}
	}
	return ""
}

// EnsureVolumes creates volumes where needed and returns a list of volumes that were created.
func EnsureVolumes(ctx context.Context, logger logrus.FieldLogger, stacks []atlasfile.StackConfig, a *atlasfile.Atlasfile) (EnsuredVolumes, error) {
	ensuredVolumes := make([]EnsuredVolume, 0)

	for _, stack := range stacks {
		for _, stackService := range stack.Services {
			service := a.GetService(stackService.Name)
			for _, volume := range service.Volumes {
				if volume.IsVolume {
					// Create volume *per stack*
					volName := helper.RandomizedName(fmt.Sprintf("atlas-%s-%s-%s", stack.Name, service.Name, volume.HostPathOrVolumeName))
					err := CreateVolume(ctx, logger, volName)
					if err != nil {
						return nil, fmt.Errorf("could not create volume: %w", err)
					}

					ensuredVolumes = append(ensuredVolumes, EnsuredVolume{
						Stack:        stack.Name,
						Service:      service.Name,
						VolumeName:   volume.HostPathOrVolumeName,
						PhysicalName: volName,
					})
				}
			}
		}
	}

	return ensuredVolumes, nil
}
