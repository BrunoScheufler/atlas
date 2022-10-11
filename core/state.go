package atlas

import (
	"encoding/json"
	"fmt"
	"github.com/brunoscheufler/atlas/atlasfile"
	"github.com/brunoscheufler/atlas/docker"
	"github.com/sirupsen/logrus"
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
	for i := range stackNames {
		for j := range s.Stacks {
			if s.Stacks[j].Name == stackNames[i] {
				stacks = append(stacks, s.Stacks[j])
				break
			}
		}

		return nil, fmt.Errorf("could not find stack %s in state file", stackNames[i])
	}

	return stacks, nil
}

type StateStack struct {
	Name     string         `json:"name"`
	Network  string         `json:"network"`
	Services []StateService `json:"services"`
}

type StateService struct {
	Name          string `json:"name"`
	ContainerName string `json:"containerName"`
}

func getStatefilePath(cwd string) string {
	return filepath.Join(cwd, ".atlas", "state.json")
}

func readState(rootDir, version string, logger logrus.FieldLogger) (*Statefile, error) {
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

func writeState(rootDir, version string, stacks []atlasfile.StackConfig, volumes docker.EnsuredVolumes) error {
	stateStacks := make([]StateStack, len(stacks))
	for i := range stacks {
		services := make([]StateService, len(stacks[i].Services))
		stack := stacks[i]
		for j := range stack.Services {
			svc := stack.Services[j]
			services[j] = StateService{
				Name:          svc.Name,
				ContainerName: stack.GetContainerName(svc.Name),
			}
		}

		stateStacks[i] = StateStack{
			Name:     stack.Name,
			Services: services,
			Network:  stack.GetNetworkName(),
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
