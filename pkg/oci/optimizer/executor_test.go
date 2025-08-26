package optimizer

import (
	"context"
	"github.com/jessesanford/idpbuilder/pkg/oci/api"
	"testing"
)

func TestNewExecutor(t *testing.T) {
	exec := NewExecutor(4)
	if exec == nil || exec.workers != 4 {
		t.Error("expected executor with 4 workers")
	}
}

func TestExecuteEmpty(t *testing.T) {
	exec := NewExecutor(2)
	result, err := exec.Execute(context.Background(), []*api.Stage{})
	if err != nil || !result.Success {
		t.Error("expected success for empty stages")
	}
}
