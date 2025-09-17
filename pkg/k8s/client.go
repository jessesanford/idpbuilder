// pkg/k8s/client.go
package k8s

import (
	"k8s.io/client-go/kubernetes"
)

// Client wraps kubernetes clientset
type Client struct {
	clientset kubernetes.Interface
}

// NewClient creates a new k8s client
func NewClient() (*Client, error) {
	return &Client{}, nil
}