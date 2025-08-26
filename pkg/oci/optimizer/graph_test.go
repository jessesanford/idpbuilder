package optimizer

import (
	"testing"

	"github.com/jessesanford/idpbuilder/pkg/oci/api"
)

func TestNewGraphBuilder(t *testing.T) {
	gb := NewGraphBuilder()
	if gb == nil {
		t.Fatal("expected graph builder instance")
	}
	if gb.nodes == nil {
		t.Error("expected initialized nodes map")
	}
}

func TestBuildEmptyStages(t *testing.T) {
	gb := NewGraphBuilder()
	
	_, err := gb.Build([]*api.Stage{})
	if err == nil {
		t.Error("expected error for empty stages")
	}
}

func TestBuildSingleStage(t *testing.T) {
	gb := NewGraphBuilder()
	stages := []*api.Stage{{Name: "test", BaseImage: "alpine"}}
	
	graph, err := gb.Build(stages)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(graph.Nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(graph.Nodes))
	}
}