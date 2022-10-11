package atlas

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/brunoscheufler/atlas/atlasfile"
	"github.com/brunoscheufler/atlas/docker"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"os"
	"path/filepath"
)

type Statefile struct {
	Version string       `json:"version"`
	Stacks  []StateStack `json:"stacks"`
	Volumes []string     `json:"volumes"`
}

func (s *Statefile) GetStacks(stackNames []string) ([]StateStack, error) {
	if len(stackNames) == 0 {
		return s.Stacks, nil
	}

	stacks := make([]StateStack, 0)
	for _, stackName := range stackNames {
		for _, stack := range s.Stacks {
			if stack.Name == stackName {
				stacks = append(stacks, stack)
				break
			}
		}
	}

	return stacks, nil
}

type StateStack struct {
	Name     string         `json:"name"`
	Network  string         `json:"network"`
	Services []StateService `json:"services"`
}

type StateService struct {
	Name string `json:"name"`

	ContainerName  string                 `json:"containerName"`
	ContainerInfos *docker.ContainerInfos `json:"containerInfo"`
}

func getStatefilePath(cwd string) string {
	return filepath.Join(cwd, ".atlas", "state.json")
}

func refreshState(ctx context.Context, rootDir string, stateFile *Statefile) error {
	newStacks := make([]StateStack, 0)

	for _, stack := range stateFile.Stacks {
		// Refresh service containers
		currentServices := make([]StateService, 0)
		for _, service := range stack.Services {
			infos, err := docker.GetContainerInfo(ctx, service.ContainerName)
			if err != nil {
				return fmt.Errorf("could not get container info: %w", err)
			}

			if infos == nil {
				continue
			}

			currentServices = append(currentServices, StateService{
				Name:           service.Name,
				ContainerName:  service.ContainerName,
				ContainerInfos: infos,
			})
		}

		stack.Services = currentServices

		// Refresh network
		networkId, err := docker.GetNetworkId(ctx, stack.Network)
		if err != nil {
			return fmt.Errorf("could not get network id: %w", err)
		}

		if networkId != "" && len(currentServices) > 0 {
			newStacks = append(newStacks, stack)
		}
	}

	stateFile.Stacks = newStacks

	err := writeStateFileRaw(rootDir, stateFile)
	if err != nil {
		return fmt.Errorf("could not write state file: %w", err)
	}

	return nil
}

func readState(ctx context.Context, rootDir, version string, logger logrus.FieldLogger) (*Statefile, error) {
	stateFile := Statefile{}

	stateFilePath := getStatefilePath(rootDir)

	_, err := os.Stat(stateFilePath)
	if err != nil {
		if os.IsNotExist(err) {

			return nil, nil
		}

		return nil, fmt.Errorf("could not stat state file: %w", err)
	}

	marshalled, err := os.ReadFile(stateFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not read state file: %w", err)
	}

	err = json.Unmarshal(marshalled, &stateFile)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal state file: %w", err)
	}

	if stateFile.Version != version {
		_ = clearStatefile(rootDir)
		logger.Warnf("Found an existing state file from older (current: %s, stored: %s) version. Please clean remaining resources manually or run atlas down --all to remove all containers and volumes.\n", version, stateFile.Version)
		return nil, nil
	}

	err = refreshState(ctx, rootDir, &stateFile)
	if err != nil {
		return nil, fmt.Errorf("could not refresh state: %w", err)
	}

	return &stateFile, nil
}

func clearStatefile(rootDir string) error {
	stateFilePath := getStatefilePath(rootDir)

	err := os.Remove(stateFilePath)
	if err != nil {
		return fmt.Errorf("could not remove state file: %w", err)
	}

	return nil
}

func writeState(ctx context.Context, rootDir, version string, stacks []atlasfile.StackConfig, volumes docker.EnsuredVolumes, networks docker.EnsuredNetworks) error {
	stateStacks := make([]StateStack, len(stacks))

	for i := range stacks {
		services := make([]StateService, len(stacks[i].Services))
		stack := stacks[i]

		{
			g, ctx := errgroup.WithContext(ctx)
			for j := range stack.Services {
				j := j
				svc := stack.Services[j]

				g.Go(func() error {
					containerName := stack.GetContainerName(svc.Name)
					containerInfos, err := docker.GetContainerInfo(ctx, containerName)
					if err != nil {
						return fmt.Errorf("could not get container infos: %w", err)
					}

					services[j] = StateService{
						Name:           svc.Name,
						ContainerName:  containerName,
						ContainerInfos: containerInfos,
					}

					return nil
				})
			}

			err := g.Wait()
			if err != nil {
				return fmt.Errorf("could not write container state: %w", err)
			}
		}

		stateStacks[i] = StateStack{
			Name:     stack.Name,
			Services: services,
			Network:  networks.Get(stack.Name),
		}
	}

	volumeNames := make([]string, len(volumes))
	for i := range volumes {
		volumeNames[i] = volumes[i].PhysicalName
	}

	stateFile := Statefile{
		Version: version,
		Stacks:  stateStacks,
		Volumes: volumeNames,
	}

	return writeStateFileRaw(rootDir, &stateFile)
}

func writeStateFileRaw(rootDir string, stateFile *Statefile) error {
	marshalled, err := json.Marshal(stateFile)
	if err != nil {
		return fmt.Errorf("could not marshal state file: %w", err)
	}

	err = os.WriteFile(getStatefilePath(rootDir), marshalled, 0644)
	if err != nil {
		return fmt.Errorf("could not write state file: %w", err)
	}

	return nil
}
