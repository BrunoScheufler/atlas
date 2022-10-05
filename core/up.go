package atlas

import (
	"context"
	"fmt"
	"github.com/brunoscheufler/atlas/atlasfile"
	"github.com/brunoscheufler/atlas/docker"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
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

	stacks, err := mergedFile.GetStacks(stackNames)
	if err != nil {
		return fmt.Errorf("could not get stacks: %w", err)
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

	err = ensureVolumes(ctx, logger, stacks, mergedFile)
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
	for j := range stack.Services {
		stackService := &stack.Services[j]
		service := file.GetService(stackService.Name)

		logger.WithField("stack", stack.Name).Infoln(fmt.Sprintf("Starting %s", services[j].Name))

		containerName := randomizedName(fmt.Sprintf("atlas-%s-%s", stack.Name, service.Name))

		err := docker.CreateServiceContainer(ctx, logger, stack, service, stackService, file, containerName)
		if err != nil {
			return fmt.Errorf("could not create service container: %w", err)
		}

		stack.SetContainerName(service.Name, containerName)
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

// ensureVolumes creates volumes where needed (volumes are created per stack and stored in the stack service)
func ensureVolumes(ctx context.Context, logger logrus.FieldLogger, stacks []atlasfile.StackConfig, a *atlasfile.Atlasfile) error {
	for i := range stacks {
		stack := &stacks[i]
		for j := range stack.Services {
			stackService := &stack.Services[j]
			service := a.GetService(stackService.Name)
			for k := range service.Volumes {
				volume := &service.Volumes[k]
				if volume.IsVolume {
					// Create volume *per stack*
					volName := randomizedName(fmt.Sprintf("atlas-%s-%s-%s", stack.Name, service.Name, volume.HostPathOrVolumeName))
					err := docker.CreateVolume(ctx, logger, volName)
					if err != nil {
						return fmt.Errorf("could not create volume: %w", err)
					}

					stackService.SetVolumeName(volume.HostPathOrVolumeName, volName)
				}
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

func buildArtifacts(ctx context.Context, logger logrus.FieldLogger, file *atlasfile.Atlasfile, layers [][]string, cwd string) error {
	for _, layer := range layers {
		g, ctx := errgroup.WithContext(ctx)

		for _, artifactName := range layer {
			artifactName := artifactName

			g.Go(func() error {
				artifact := file.GetArtifact(artifactName)
				if artifact == nil {
					return fmt.Errorf("could not find artifact %s", artifactName)
				}

				err := docker.BuildArtifact(ctx, logger, artifact, cwd)
				if err != nil {
					return fmt.Errorf("could not build artifact %s: %w", artifact.Name, err)
				}

				return nil
			})
		}

		err := g.Wait()
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO support caching -> only build when artifact inputs changed
