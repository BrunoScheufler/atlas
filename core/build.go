package atlas

import (
	"context"
	"fmt"
	"github.com/brunoscheufler/atlas/atlasfile"
	"github.com/brunoscheufler/atlas/docker"
	"github.com/sirupsen/logrus"
)

func Build(ctx context.Context, logger logrus.FieldLogger, version, cwd string, stackNames []string) error {
	logger.WithFields(
		logrus.Fields{
			"version": version,
			"cwd":     cwd,
			"stacks":  stackNames,
		},
	).Debugf("Running core.Build")

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

	err = buildArtifacts(ctx, logger, mergedFile, layers)
	if err != nil {
		return fmt.Errorf("could not build artifacts: %w", err)
	}

	return nil
}
