package extractor

import (
	"bytes"
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

// GetGiteaPod finds the Gitea pod in the Kind cluster
func (e *kindExtractor) GetGiteaPod(ctx context.Context) (podName, namespace string, err error) {
	// Common namespaces where Gitea might be installed
	namespaces := []string{"gitea", "default", "gitea-system"}

	for _, ns := range namespaces {
		// Try to find Gitea pods using label selectors
		pods, err := e.client.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{
			LabelSelector: "app=gitea",
		})

		if err != nil {
			continue // Try next namespace
		}

		// Find a running pod
		for _, pod := range pods.Items {
			if pod.Status.Phase == corev1.PodRunning {
				return pod.Name, pod.Namespace, nil
			}
		}

		// Also try app.kubernetes.io/name label
		pods, err = e.client.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{
			LabelSelector: "app.kubernetes.io/name=gitea",
		})

		if err != nil {
			continue
		}

		for _, pod := range pods.Items {
			if pod.Status.Phase == corev1.PodRunning {
				return pod.Name, pod.Namespace, nil
			}
		}
	}

	return "", "", &PodNotFoundError{
		AppName:    "gitea",
		Namespaces: namespaces,
	}
}

// execInPod executes a command in a pod and returns the output
func (e *kindExtractor) execInPod(ctx context.Context, namespace, podName string, command []string) ([]byte, error) {
	// Create exec request
	req := e.client.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command: command,
			Stdout:  true,
			Stderr:  true,
			TTY:     false,
		}, scheme.ParameterCodec)

	// Get the config for the exec
	config, err := clientcmd.BuildConfigFromFlags("", e.kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	// Create the executor
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return nil, fmt.Errorf("failed to create executor: %w", err)
	}

	// Prepare stdout and stderr buffers
	var stdout, stderr bytes.Buffer

	// Execute the command
	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	})

	if err != nil {
		return nil, &CertificateReadError{
			PodName:   podName,
			Namespace: namespace,
			Path:      fmt.Sprintf("command: %v", command),
			Err:       fmt.Errorf("exec failed: %w, stderr: %s", err, stderr.String()),
		}
	}

	return stdout.Bytes(), nil
}