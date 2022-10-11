package docker

import (
	"context"
	"fmt"
	"github.com/brunoscheufler/atlas/atlasfile"
	"github.com/brunoscheufler/atlas/exec"
	"github.com/brunoscheufler/atlas/helper"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
	"time"
)

func CreateServiceContainer(
	ctx context.Context,
	logger logrus.FieldLogger,
	stack *atlasfile.StackConfig,
	service *atlasfile.ServiceConfig,
	stackService *atlasfile.StackService,
	file *atlasfile.Atlasfile,
	ensuredVolumes EnsuredVolumes,
	ensuredNetworks EnsuredNetworks,
	containerName string,
) error {
	args := []string{
		"run",
		"-d",
		"--name",
		containerName,
		"--hostname",
		service.Name,
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
			volName := volume.GetVolumeNameOrHostPath(filepath.Dir(service.GetDirpath()), ensuredVolumes.Get(stack.Name, service.Name, volume.HostPathOrVolumeName))
			args = append(args, "-v", fmt.Sprintf("%s:%s", volName, volume.ContainerPath))
		}
	}

	netName := ensuredNetworks.Get(stack.Name)
	if netName != "" {
		args = append(args, "--network", netName)
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

	if service.Entrypoint != nil {
		args = append(args, "--entrypoint", strings.Join(service.Entrypoint, " "))
	}

	if service.Interactive {
		args = append(args, "-i")
	}

	if service.TTY {
		args = append(args, "-t")
	}

	imageName, err := file.GetServiceImage(service)
	if err != nil {
		return fmt.Errorf("could not get service image: %w", err)
	}
	args = append(args, imageName)

	if service.Command != nil {
		args = append(args, service.Command...)
	}

	err = exec.RunCommand(ctx, logger, fmt.Sprintf("docker %s", strings.Join(args, " ")), exec.RunCommandOptions{})
	if err != nil {
		return fmt.Errorf("could not create container %s: %w", containerName, err)
	}

	if stackService.JoinStackNetworks != nil {
		for _, stackName := range stackService.JoinStackNetworks {
			netName := ensuredNetworks.Get(stackName)
			if netName == "" {
				return fmt.Errorf("could not find network for stack %s", stackName)
			}

			err = exec.RunCommand(ctx, logger, fmt.Sprintf("docker network connect %s %s", netName, containerName), exec.RunCommandOptions{})
			if err != nil {
				return fmt.Errorf("could not connect container %s to network %s: %w", containerName, netName, err)
			}
		}
	}

	return nil
}

type ContainerInfos struct {
	FetchedAt string `json:"fetchedAt"`
	Id        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	State     string `json:"state"`
}

func GetContainerInfo(ctx context.Context, containerName string) (*ContainerInfos, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("could not create docker client: %w", err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		All: true,
		Filters: filters.NewArgs(
			filters.Arg("name", containerName),
		),
	})
	if err != nil {
		return nil, fmt.Errorf("could not list containers: %w", err)
	}

	if len(containers) == 0 {
		return nil, nil
	}

	return &ContainerInfos{
		FetchedAt: time.Now().Format(time.RFC3339),
		Id:        containers[0].ID,
		Name:      containers[0].Names[0],
		Status:    containers[0].Status,
		State:     containers[0].State,
	}, nil
}

func StartContainer(ctx context.Context, containerName string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("could not create docker client: %w", err)
	}

	err = cli.ContainerStart(ctx, containerName, types.ContainerStartOptions{})
	if err != nil {
		return fmt.Errorf("could not start container %s: %w", containerName, err)
	}

	return nil
}

func StopContainer(ctx context.Context, containerName string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("could not create docker client: %w", err)
	}

	err = cli.ContainerStop(ctx, containerName, nil)
	if err != nil {
		return fmt.Errorf("could not stop container %s: %w", containerName, err)
	}

	return nil
}
