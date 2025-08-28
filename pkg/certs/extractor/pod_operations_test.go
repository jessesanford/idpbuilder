package extractor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetGiteaPod_Success(t *testing.T) {
	// Create fake kubernetes client
	fakeClient := fake.NewSimpleClientset()

	// Create a running Gitea pod
	giteaPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gitea-12345",
			Namespace: "gitea",
			Labels: map[string]string{
				"app": "gitea",
			},
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
		},
	}
	fakeClient.CoreV1().Pods("gitea").Create(context.TODO(), giteaPod, metav1.CreateOptions{})

	extractor := &kindExtractor{
		client: fakeClient,
	}

	// Test finding the pod
	podName, namespace, err := extractor.GetGiteaPod(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "gitea-12345", podName)
	assert.Equal(t, "gitea", namespace)
}

func TestGetGiteaPod_AlternativeLabel(t *testing.T) {
	// Create fake kubernetes client
	fakeClient := fake.NewSimpleClientset()

	// Create a running Gitea pod with alternative label
	giteaPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gitea-server-abcde",
			Namespace: "default",
			Labels: map[string]string{
				"app.kubernetes.io/name": "gitea",
			},
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
		},
	}
	fakeClient.CoreV1().Pods("default").Create(context.TODO(), giteaPod, metav1.CreateOptions{})

	extractor := &kindExtractor{
		client: fakeClient,
	}

	// Test finding the pod with alternative label
	podName, namespace, err := extractor.GetGiteaPod(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "gitea-server-abcde", podName)
	assert.Equal(t, "default", namespace)
}

func TestGetGiteaPod_NotFound(t *testing.T) {
	// Create fake kubernetes client with no pods
	fakeClient := fake.NewSimpleClientset()

	extractor := &kindExtractor{
		client: fakeClient,
	}

	// Test when no Gitea pod is found
	podName, namespace, err := extractor.GetGiteaPod(context.Background())
	assert.Error(t, err)
	assert.Empty(t, podName)
	assert.Empty(t, namespace)

	var podErr *PodNotFoundError
	assert.ErrorAs(t, err, &podErr)
	assert.Equal(t, "gitea", podErr.AppName)
	assert.Contains(t, podErr.Namespaces, "gitea")
	assert.Contains(t, podErr.Namespaces, "default")
	assert.Contains(t, podErr.Namespaces, "gitea-system")
}

func TestGetGiteaPod_MultiplePodsSelectRunning(t *testing.T) {
	// Create fake kubernetes client
	fakeClient := fake.NewSimpleClientset()

	// Create a pod that's not running
	pendingPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gitea-pending",
			Namespace: "gitea",
			Labels: map[string]string{
				"app": "gitea",
			},
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodPending,
		},
	}
	fakeClient.CoreV1().Pods("gitea").Create(context.TODO(), pendingPod, metav1.CreateOptions{})

	// Create a running pod
	runningPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gitea-running",
			Namespace: "gitea",
			Labels: map[string]string{
				"app": "gitea",
			},
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
		},
	}
	fakeClient.CoreV1().Pods("gitea").Create(context.TODO(), runningPod, metav1.CreateOptions{})

	extractor := &kindExtractor{
		client: fakeClient,
	}

	// Test that running pod is selected
	podName, namespace, err := extractor.GetGiteaPod(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "gitea-running", podName)
	assert.Equal(t, "gitea", namespace)
}

func TestPodNotFoundError(t *testing.T) {
	err := &PodNotFoundError{
		AppName:    "gitea",
		Namespaces: []string{"gitea", "default", "gitea-system"},
	}

	errMsg := err.Error()
	assert.Contains(t, errMsg, "no running gitea pod found")
	assert.Contains(t, errMsg, "[gitea default gitea-system]")
	assert.Contains(t, errMsg, "ensure gitea is installed and running")
}

func TestCertificateReadError(t *testing.T) {
	originalErr := assert.AnError
	err := &CertificateReadError{
		PodName:   "test-pod",
		Namespace: "test-ns",
		Path:      "/path/to/cert",
		Err:       originalErr,
	}

	errMsg := err.Error()
	assert.Contains(t, errMsg, "failed to read certificate from test-ns/test-pod")
	assert.Contains(t, errMsg, "/path/to/cert")

	// Test Unwrap
	assert.Equal(t, originalErr, err.Unwrap())
}