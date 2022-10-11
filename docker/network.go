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
)

func CreateNetwork(ctx context.Context, logger logrus.FieldLogger, name string) error {
	err := exec.RunCommand(ctx, logger, fmt.Sprintf("docker network create %s", name), exec.RunCommandOptions{})
	if err != nil {
		return fmt.Errorf("could not create network %s: %w", name, err)
	}

	return nil
}

func GetNetworkId(ctx context.Context, networkName string) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", fmt.Errorf("could not create docker client: %w", err)
	}

	networks, err := cli.NetworkList(ctx, types.NetworkListOptions{
		Filters: filters.NewArgs(
			filters.Arg("name", networkName),
		),
	})
	if err != nil {
		return "", fmt.Errorf("could not list networks: %w", err)
	}

	if len(networks) == 0 {
		return "", nil
	}

	return networks[0].ID, nil
}

type EnsuredNetwork struct {
	Stack        string
	PhysicalName string
}

type EnsuredNetworks []EnsuredNetwork

func (e *EnsuredNetworks) Get(stackName string) string {
	for _, net := range *e {
		if net.Stack == stackName {
			return net.PhysicalName
		}
	}
	return ""
}

// EnsureNetworks creates networks where needed and returns a list of networks that were created.
func EnsureNetworks(ctx context.Context, logger logrus.FieldLogger, stacks []atlasfile.StackConfig, a *atlasfile.Atlasfile) (EnsuredNetworks, error) {
	ensuredNetworks := make([]EnsuredNetwork, 0)

	for _, stack := range stacks {
		netName := helper.RandomizedName(fmt.Sprintf("atlas-%s", stack.Name))

		err := CreateNetwork(ctx, logger, netName)
		if err != nil {
			return nil, fmt.Errorf("could not create network: %w", err)
		}

		ensuredNetworks = append(ensuredNetworks, EnsuredNetwork{
			Stack:        stack.Name,
			PhysicalName: netName,
		})
	}

	return ensuredNetworks, nil
}
