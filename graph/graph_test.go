package graph

import (
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTopologicalSortWithLayers(t *testing.T) {
	g := New[string]()

	g.AddNode("Frontend URL")
	g.AddNode("API URL")
	g.AddNode("Auth URL")

	g.AddNode("Queue")
	g.AddNode("S3 Bucket")
	g.AddNode("SES Mailing")

	g.AddNode("Database Instance")

	g.AddNode("Seed Migration")
	g.AddEdge("Database Instance", "Seed Migration")

	g.AddNode("Database")
	g.AddEdge("Database Instance", "Database")
	g.AddEdge("Seed Migration", "Database")

	g.AddNode("Frontend Deployment")
	g.AddNode("Frontend DNS Record")

	g.AddNode("Load Balancer")
	g.AddNode("API DNS Record")
	g.AddNode("Auth DNS Record")

	g.AddNode("API Deployment")
	g.AddEdge("Load Balancer", "API Deployment")
	g.AddEdge("Database", "API Deployment")
	g.AddEdge("Queue", "API Deployment")
	g.AddEdge("S3 Bucket", "API Deployment")
	g.AddEdge("SES Mailing", "API Deployment")
	g.AddEdge("API DNS Record", "API Deployment")

	g.AddNode("Auth Deployment")
	g.AddEdge("Load Balancer", "Auth Deployment")
	g.AddEdge("Database", "Auth Deployment")
	g.AddEdge("Auth DNS Record", "Auth Deployment")

	g.AddNode("Workflows Deployment")
	g.AddEdge("Database", "Workflows Deployment")
	g.AddEdge("Queue", "Workflows Deployment")
	g.AddEdge("S3 Bucket", "Workflows Deployment")
	g.AddEdge("SES Mailing", "Workflows Deployment")

	g.AddEdge("Frontend URL", "Frontend DNS Record")
	g.AddEdge("Frontend Deployment", "Frontend DNS Record")

	g.AddEdge("Frontend URL", "Frontend Deployment")
	g.AddEdge("API URL", "Frontend Deployment")

	g.AddEdge("Load Balancer", "API DNS Record")
	g.AddEdge("Load Balancer", "Auth DNS Record")

	resp, err := g.TopologicalSortWithLayers()
	assert.NoError(t, err)

	cupaloy.SnapshotT(t, resp)
}

func TestDFS(t *testing.T) {
	g := New[string]()
	g.AddNode("a")
	g.AddNode("b")
	g.AddNode("c")
	g.AddNode("d")
	g.AddNode("e")

	g.AddEdge("a", "b")
	g.AddEdge("a", "c")

	g.AddEdge("c", "d")
	g.AddEdge("b", "d")

	g.AddEdge("d", "e")

	assert.Equal(t, []string{"c", "d", "e", "b"}, g.DFS("a", -1))
}

func TestTransitiveReduction(t *testing.T) {
	g := New[string]()

	g.AddNode("a")
	g.AddNode("b")
	g.AddNode("c")
	g.AddNode("d")
	g.AddNode("e")

	g.AddEdge("a", "b")
	g.AddEdge("a", "c")
	g.AddEdge("a", "d")

	g.AddEdge("b", "d")
	g.AddEdge("c", "d")

	g.AddEdge("c", "e")
	g.AddEdge("d", "e")

	reduced, err := g.TransitiveReduction()
	if err != nil {
		t.Fatalf("could not reduce graph: %s", err.Error())
	}

	assert.Equal(t, []string{"a", "b", "c", "d", "e"}, reduced.Nodes())
	assert.Equal(t, [][]string{
		{"a", "b"},
		{"a", "c"},
		{"b", "d"},
		{"c", "d"},
		{"d", "e"},
	}, reduced.Edges())
}
