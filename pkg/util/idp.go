package util

import (
	"context"
	"fmt"
	"github.com/cnoe-io/idpbuilder/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetConfig(ctx context.Context) (v1alpha1.BuildCustomizationSpec, error) {
	b := v1alpha1.BuildCustomizationSpec{}

	kubeConfig, err := GetKubeConfig()
	if err != nil {
		return b, fmt.Errorf("getting kube config: %w", err)
	}

	kubeClient, err := GetKubeClient(kubeConfig)
	if err != nil {
		return b, fmt.Errorf("getting kube client: %w", err)
	}

	list, err := getLocalBuild(ctx, kubeClient)
	if err != nil {
		return b, err
	}

	// Phase 1 Design Decision: Single LocalBuild Assumption
	// GetLocalBuild returns the first LocalBuild resource found.
	// Phase 1 deployment model supports one LocalBuild per cluster.
	// Multi-LocalBuild support is deferred to Phase 2+ as needed.
	if len(list.Items) == 0 {
		return b, fmt.Errorf("no LocalBuild resource found")
	}

	return list.Items[0].Spec.BuildCustomization, nil
}

func getLocalBuild(ctx context.Context, kubeClient client.Client) (v1alpha1.LocalbuildList, error) {
	localBuildList := v1alpha1.LocalbuildList{}
	return localBuildList, kubeClient.List(ctx, &localBuildList)
}
