package optimizer

import (
	"context"
	"testing"
	"time"

	"github.com/jessesanford/idpbuilder/pkg/oci/api"
)

func TestNewExecutor(t *testing.T) {
	exec := NewExecutor(4)
	if exec == nil {
		t.Fatal("expected executor instance")
	}
	if exec.workers != 4 {
		t.Errorf("expected 4 workers, got %d", exec.workers)
	}
}

func TestExecuteEmptyStages(t *testing.T) {
	exec := NewExecutor(2)
	ctx := context.Background()
	
	result, err := exec.Execute(ctx, []*api.Stage{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Error("expected success for empty stages")
	}
}