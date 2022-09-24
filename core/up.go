package atlas

import (
	"context"
	"fmt"
	"github.com/brunoscheufler/atlas/atlasfile"
	"github.com/brunoscheufler/atlas/docker"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

func Up(ctx context.Context, logger logrus.FieldLogger, version, cwd string, stackNames []string) error {
	logger.WithFields(
		logrus.Fields{
			"version": version,
			"cwd":     cwd,
			"stacks":  stackNames,
		},
	).Debugf("Running core.Up")

	cwd, err := atlasfile.FindRootDir(cwd)
	if err != nil {
		return fmt.Errorf("could not find root directory: %w", err)
	}

	logger.WithField("cwd", cwd).Debugf("Found root directory")

	mergedFile, err := atlasfile.Collect(ctx, logger, version, cwd)
	if err != nil {
		return fmt.Errorf("could not collect atlas files: %w", err)
	}

	if !docker.IsRunning(ctx) {
		return fmt.Errorf("docker is not running")
	}

	err = Down(ctx, logger, cwd, stackNames)
	if err != nil {
		return fmt.Errorf("could not down: %w", err)
	}

	var stacks []atlasfile.StackConfig
	if stackNames != nil {
		for _, name := range stackNames {
			stack := mergedFile.GetStack(name)
			if stack == nil {
				return fmt.Errorf("could not find stack %s", name)
			}
			stacks = append(stacks, *stack)
		}
	} else {
		stacks = mergedFile.Stacks
	}

	services, err := getRequiredServicesForStacks(stacks, mergedFile)
	if err != nil {
		return fmt.Errorf("could not get required services: %w", err)
	}

	immediateArtifacts, err := getImmediateArtifactsNeededByServices(services, mergedFile)

	// Build artifacts
	artifactGraph, err := buildArtifactGraphWithImmediate(mergedFile, immediateArtifacts)

	layers, err := artifactGraph.TopologicalSortWithLayers()
	if err != nil {
		return fmt.Errorf("could not topologically sort artifacts: %w", err)
	}

	// TODO only build artifacts required for stacks
	err = buildArtifacts(ctx, logger, mergedFile, layers, cwd)
	if err != nil {
		return fmt.Errorf("could not build artifacts: %w", err)
	}

	for i, stack := range stacks {
		netName := randomizedName(fmt.Sprintf("atlas-%s", stack.Name))
		stacks[i].SetNetworkName(netName)

		err = docker.CreateNetwork(ctx, logger, netName)
		if err != nil {
			return fmt.Errorf("could not create network: %w", err)
		}
	}

	err = ensureVolumes(ctx, logger, services)
	if err != nil {
		return fmt.Errorf("could not ensure volumes: %w", err)
	}

	for i := range stacks {
		logger.Infof("Launching stack %s\n", stacks[i].Name)

		err := startStack(ctx, logger, &stacks[i], mergedFile, services)
		if err != nil {
			return fmt.Errorf("could not start stack %q: %w", stacks[i].Name, err)
		}
	}

	// TODO persist session for subsequent commands

	return nil
}

func startStack(ctx context.Context, logger logrus.FieldLogger, stack *atlasfile.StackConfig, file *atlasfile.Atlasfile, services []atlasfile.ServiceConfig) error {
	bar := progressbar.NewOptions(len(stack.Services), progressbar.OptionSetDescription("Starting services"), progressbar.OptionClearOnFinish())
	defer func() {
		_ = bar.Finish()
		_ = bar.Clear()
		_ = bar.Close()
	}()

	for j := range stack.Services {
		stackService := &stack.Services[j]
		service := file.GetService(stackService.Name)

		bar.Describe(fmt.Sprintf("Starting %s", services[j].Name))

		containerName := randomizedName(fmt.Sprintf("atlas-%s-%s", stack.Name, service.Name))

		err := docker.CreateServiceContainer(ctx, logger, stack, service, stackService, file, containerName)
		if err != nil {
			return fmt.Errorf("could not create service container: %w", err)
		}

		stack.SetContainerName(service.Name, containerName)

		_ = bar.Add(1)
	}

	return nil
}

func getImmediateArtifactsNeededByServices(services []atlasfile.ServiceConfig, file *atlasfile.Atlasfile) ([]atlasfile.ArtifactConfig, error) {
	var artifacts []atlasfile.ArtifactConfig

	for _, service := range services {
		if service.Artifact == nil {
			continue
		}

		var artifactName string
		if service.Artifact.Name != "" {
			artifactName = service.Artifact.Name
		} else {
			artifactName = service.Artifact.Artifact.Name

		}

		artifact := file.GetArtifact(artifactName)
		if artifact == nil {
			return nil, fmt.Errorf("could not find artifact %s", artifactName)
		}

		artifacts = append(artifacts, *artifact)
	}

	return artifacts, nil
}

func ensureVolumes(ctx context.Context, logger logrus.FieldLogger, services []atlasfile.ServiceConfig) error {
	for i, service := range services {
		for j, volume := range service.Volumes {
			if volume.IsVolume {
				volName := randomizedName(fmt.Sprintf("atlas-%s-%s", service.Name, volume.HostPathOrVolumeName))
				err := docker.CreateVolume(ctx, logger, volName)
				if err != nil {
					return fmt.Errorf("could not create volume: %w", err)
				}

				services[i].Volumes[j].SetVolName(volName)
			}
		}
	}
	return nil
}

func randomizedName(name string) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	digits := 4
	var suffix string
	for i := 0; i < digits; i++ {
		suffix += string(rune(rnd.Intn(10) + 48))
	}

	return fmt.Sprintf("%s-%s", name, suffix)
}

func getRequiredServicesForStacks(stacks []atlasfile.StackConfig, file *atlasfile.Atlasfile) ([]atlasfile.ServiceConfig, error) {
	var services []atlasfile.ServiceConfig

	for _, stack := range stacks {
		services = append(services, getRequiredServices(stack, file)...)
	}

	return services, nil
}

func getRequiredServices(stack atlasfile.StackConfig, file *atlasfile.Atlasfile) []atlasfile.ServiceConfig {
	services := make([]atlasfile.ServiceConfig, len(stack.Services))
	for i2, stackService := range stack.Services {
		for _, service := range file.Services {
			if service.Name == stackService.Name {
				services[i2] = service
			}
		}
	}

	return services
}

func buildArtifacts(ctx context.Context, logger logrus.FieldLogger, file *atlasfile.Atlasfile, layers [][]string, rootDir string) error {
	for _, layer := range layers {
		bar := progressbar.NewOptions(len(layer), progressbar.OptionSetDescription("Building artifacts"), progressbar.OptionClearOnFinish())

		// TODO Run in parallel
		for _, artifactName := range layer {
			artifact := file.GetArtifact(artifactName)
			if artifact == nil {
				return fmt.Errorf("could not find artifact %s", artifactName)
			}

			bar.Describe(fmt.Sprintf("Building artifact %s", artifactName))

			err := docker.BuildArtifact(ctx, logger, artifact, rootDir)
			if err != nil {
				return fmt.Errorf("could not build artifact %s: %w", artifact.Name, err)
			}

			_ = bar.Add(1)
		}

		_ = bar.Clear()
		_ = bar.Close()
	}

	return nil
}

// TODO support caching -> only build when artifact inputs changed
