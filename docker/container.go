package docker

import (
	"context"
	"fmt"
	"github.com/brunoscheufler/atlas/atlasfile"
	"github.com/brunoscheufler/atlas/exec"
	"github.com/brunoscheufler/atlas/helper"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
)

func CreateServiceContainer(
	ctx context.Context,
	logger logrus.FieldLogger,
	stack *atlasfile.StackConfig,
	service *atlasfile.ServiceConfig,
	stackService *atlasfile.StackService,
	file *atlasfile.Atlasfile,
	containerName string,
) error {
	args := []string{
		"run",
		"-d",
		"--name",
		containerName,
	}

	if string(service.Restart) == "" {
		service.Restart = atlasfile.ContainerRestartsAlways
	}
	args = append(args, "--restart", string(service.Restart))

	envVars := make(map[string]string)

	if service.EnvironmentFiles != nil {
		for _, file := range service.EnvironmentFiles {
			filePath := filepath.Join(filepath.Dir(service.GetDirpath()), file)

			readVars, err := helper.ReadEnvFile(filePath)
			if err != nil {
				return fmt.Errorf("could not read environment file %s: %w", file, err)
			}

			for k, v := range readVars {
				envVars[k] = v
			}
		}
	}

	if service.Environment != nil {
		for key, value := range service.Environment {
			envVars[key] = value
		}
	}

	if stackService.Environment != nil {
		for key, value := range stackService.Environment {
			envVars[key] = value
		}
	}

	for key, value := range envVars {
		args = append(args, "-e", fmt.Sprintf("%s=%q", key, value))
	}

	if service.Volumes != nil {
		for _, volume := range service.Volumes {
			volName := volume.GetVolumeNameOrHostPath(filepath.Dir(service.GetDirpath()))
			args = append(args, "-v", fmt.Sprintf("%s:%s", volName, volume.ContainerPath))
		}
	}

	if stack.GetNetworkName() != "" {
		args = append(args, "--network", stack.GetNetworkName())
	}

	if stackService.ExposePorts != nil {
		for _, expose := range stackService.ExposePorts {
			servicePortRequest := atlasfile.GetServicePort(service.Ports, expose.ContainerPort)
			if servicePortRequest == nil {
				return fmt.Errorf("could not find port %d in service %s", expose.ContainerPort, service.Name)
			}

			args = append(args, "-p", fmt.Sprintf("%d:%d/%s", expose.HostPort, expose.ContainerPort, servicePortRequest.Protocol))
		}
	}

	imageName, err := file.GetServiceImage(service)
	if err != nil {
		return fmt.Errorf("could not get service image: %w", err)
	}
	args = append(args, imageName)

	if service.Command != nil {
		args = append(args, service.Command...)
	}

	err = exec.RunCommand(ctx, logger, fmt.Sprintf("docker %s", strings.Join(args, " ")), "", nil, false)
	if err != nil {
		return fmt.Errorf("could not create container %s: %w", containerName, err)
	}

	if stackService.JoinStackNetworks != nil {
		for _, stackName := range stackService.JoinStackNetworks {
			netName := file.GetStack(stackName).GetNetworkName()

			err = exec.RunCommand(ctx, logger, fmt.Sprintf("docker network connect %s %s", netName, containerName), "", nil, false)
			if err != nil {
				return fmt.Errorf("could not connect container %s to network %s: %w", containerName, netName, err)
			}
		}
	}

	return nil
}
