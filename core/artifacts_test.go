package atlas

import (
	"github.com/bradleyjkemp/cupaloy"
	"testing"
)

func TestBuildArtifactGraph(t *testing.T) {
	g, err := buildArtifactGraph(&Atlasfile{
		Artifacts: []ArtifactConfig{
			{
				Name:      "base",
				DependsOn: ArtifactDependsOn{},
			},
			{
				Name: "api",
				DependsOn: ArtifactDependsOn{
					Artifacts: []string{"base"},
				},
			},
		},
		Services: []ServiceConfig{
			{
				Name: "api",
				Artifact: &ArtifactRef{
					Name: "api",
				},
			},
			{
				Name: "db",
				Artifact: &ArtifactRef{
					Artifact: &ArtifactConfig{
						Name: "db",
						DependsOn: ArtifactDependsOn{
							Artifacts: []string{"base"},
						},
					},
				},
			},
			{
				Name: "tool",
				Artifact: &ArtifactRef{
					Artifact: &ArtifactConfig{
						Name: "tool",
						DependsOn: ArtifactDependsOn{
							Services: []string{"api"},
						},
					},
				},
			},
		},
		Stacks: nil,
	})
	if err != nil {
		t.Fatal(err)
	}

	layers, err := g.TopologicalSortWithLayers()
	if err != nil {
		t.Fatal(err)
	}

	cupaloy.New(cupaloy.SnapshotFileExtension(".graph")).SnapshotT(t, g.String())
	cupaloy.New(cupaloy.SnapshotFileExtension(".topsort")).SnapshotT(t, layers)
}

func TestBuildArtifactGraphWithImmediate(t *testing.T) {
	testFile := mergeAtlasFiles([]Atlasfile{
		{
			Artifacts: []ArtifactConfig{
				{
					Name:      "base",
					DependsOn: ArtifactDependsOn{},
				},
				{
					Name: "api",
					DependsOn: ArtifactDependsOn{
						Artifacts: []string{"base"},
					},
				},
			},
			Services: []ServiceConfig{
				{
					Name: "api",
					Artifact: &ArtifactRef{
						Name: "api",
					},
				},
				{
					Name: "db",
					Artifact: &ArtifactRef{
						Artifact: &ArtifactConfig{
							Name: "db",
							DependsOn: ArtifactDependsOn{
								Artifacts: []string{"base"},
							},
						},
					},
				},
				{
					Name: "tool",
					Artifact: &ArtifactRef{
						Artifact: &ArtifactConfig{
							Name: "tool",
							DependsOn: ArtifactDependsOn{
								Services: []string{"api"},
							},
						},
					},
				},
			},
			Stacks: nil,
		},
	})

	immediate, err := getImmediateArtifactsNeededByServices(testFile.Services, testFile)
	if err != nil {
		t.Fatal(err)
	}

	g, err := buildArtifactGraphWithImmediate(testFile, immediate)
	if err != nil {
		t.Fatal(err)
	}

	layers, err := g.TopologicalSortWithLayers()
	if err != nil {
		t.Fatal(err)
	}

	cupaloy.New(cupaloy.SnapshotFileExtension(".graph")).SnapshotT(t, g.String())
	cupaloy.New(cupaloy.SnapshotFileExtension(".topsort")).SnapshotT(t, layers)
}
