package atlas

import (
	"fmt"
	"github.com/brunoscheufler/atlas/atlasfile"
	"github.com/brunoscheufler/atlas/graph"
)

// buildArtifactGraph builds graph of artifacts by walking all services and artifacts
func buildArtifactGraph(mergedAtlasFile *atlasfile.Atlasfile) (*graph.Graph[string], error) {
	// Build dependency graph
	artifactGraph := graph.New[string]()
	visitedServices := make(map[string]struct{})

	for _, service := range mergedAtlasFile.Services {
		err := walkService(visitedServices, mergedAtlasFile, artifactGraph, service)
		if err != nil {
			return nil, fmt.Errorf("could not walk service %s: %w", service.Name, err)
		}
	}

	err := walkArtifacts(visitedServices, mergedAtlasFile, artifactGraph, mergedAtlasFile.Artifacts)
	if err != nil {
		return nil, fmt.Errorf("could not walk artifacts: %w", err)
	}

	return artifactGraph, nil
}

// buildArtifactGraphWithImmediate builds graph of artifacts related to supplied immediateArtifacts
func buildArtifactGraphWithImmediate(mergedAtlasFile *atlasfile.Atlasfile, immediateArtifacts []atlasfile.ArtifactConfig) (*graph.Graph[string], error) {
	// Build dependency graph
	artifactGraph := graph.New[string]()
	visitedServices := make(map[string]struct{})

	err := walkArtifacts(visitedServices, mergedAtlasFile, artifactGraph, immediateArtifacts)
	if err != nil {
		return nil, fmt.Errorf("could not walk artifacts: %w", err)
	}

	return artifactGraph, nil
}

func walkArtifact(visitedServices map[string]struct{}, mergedAtlasFile *atlasfile.Atlasfile, artifactGraph *graph.Graph[string], artifact atlasfile.ArtifactConfig) error {
	if artifactGraph.HasNode(artifact.Name) {
		return nil
	}

	artifactGraph.AddNode(artifact.Name)

	// Add immediate artifact dependencies
	for _, dependsOnArtifact := range artifact.DependsOn.Artifacts {
		artifactGraph.AddNode(dependsOnArtifact)

		if !artifactGraph.HasEdge(dependsOnArtifact, artifact.Name) {
			artifactGraph.AddEdge(dependsOnArtifact, artifact.Name)
		}
	}

	// Add dependencies via services
	for _, service := range artifact.DependsOn.Services {
		svc := mergedAtlasFile.GetService(service)
		if svc == nil {
			return fmt.Errorf("service %s not found", service)
		}

		if svc.Artifact != nil {
			if svc.Artifact.Name != "" {
				artifactGraph.AddEdge(svc.Artifact.Name, artifact.Name)
			} else {
				artifactGraph.AddEdge(svc.Artifact.Artifact.Name, artifact.Name)
			}
		}

		err := walkService(visitedServices, mergedAtlasFile, artifactGraph, *svc)
		if err != nil {
			return fmt.Errorf("could not walk service %s: %w", service, err)
		}
	}

	return nil
}

func walkArtifacts(visitedServices map[string]struct{}, mergedAtlasFile *atlasfile.Atlasfile, artifactGraph *graph.Graph[string], artifacts []atlasfile.ArtifactConfig) error {
	for _, artifact := range artifacts {
		err := walkArtifact(visitedServices, mergedAtlasFile, artifactGraph, artifact)
		if err != nil {
			return fmt.Errorf("could not walk artifact %s: %w", artifact.Name, err)
		}
	}

	return nil
}

func walkService(visitedServices map[string]struct{}, mergedAtlasFile *atlasfile.Atlasfile, artifactGraph *graph.Graph[string], config atlasfile.ServiceConfig) error {
	if _, ok := visitedServices[config.Name]; ok {
		return nil
	}

	visitedServices[config.Name] = struct{}{}

	if config.Artifact == nil {
		return nil
	}

	if config.Artifact.Name != "" {
		svcArtifact := mergedAtlasFile.GetArtifact(config.Artifact.Name)
		if svcArtifact == nil {
			return fmt.Errorf("artifact %s not found", config.Artifact.Name)
		}

		err := walkArtifact(visitedServices, mergedAtlasFile, artifactGraph, *svcArtifact)
		if err != nil {
			return fmt.Errorf("could not walk artifact %s: %w", config.Artifact.Name, err)
		}

		return nil
	}

	if config.Artifact.Artifact == nil {
		return fmt.Errorf("service %s has no artifact", config.Name)
	}

	err := walkArtifact(visitedServices, mergedAtlasFile, artifactGraph, *config.Artifact.Artifact)
	if err != nil {
		return fmt.Errorf("could not walk implicit service artifact: %w", err)
	}

	return nil
}
