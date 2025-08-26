package optimizer

import (
	"github.com/jessesanford/idpbuilder/pkg/oci/api"
	"testing"
)

func TestNewGraphBuilder(t *testing.T) {
	gb := NewGraphBuilder()
	if gb == nil || gb.nodes == nil {
		t.Error("expected graph builder with initialized nodes")
	}
}

func TestBuildEmpty(t *testing.T) {
	gb := NewGraphBuilder()
	if _, err := gb.Build([]*api.Stage{}); err == nil {
		t.Error("expected error for empty stages")
	}
}

func TestBuildSingle(t *testing.T) {
	gb := NewGraphBuilder()
	stages := []*api.Stage{{Name: "test", BaseImage: "alpine"}}
	graph, err := gb.Build(stages)
	if err != nil || len(graph.Nodes) != 1 {
		t.Error("expected single node graph")
	}
}
