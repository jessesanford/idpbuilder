package registry

import (
	"time"
)

// Credentials represents authentication credentials for Gitea registry
type Credentials struct {
	Username string
	Password string
	Token    string
}

// PushOptions contains options for pushing images to registry
type PushOptions struct {
	ImageID    string
	Repository string
	Tag        string
	Insecure   bool
	Username   string
	Password   string
}

// PushResult contains the result of a push operation
type PushResult struct {
	Digest     string
	Size       int64
	PushTime   time.Duration
	Repository string
	Tag        string
}

// PullResult contains the result of a pull operation
type PullResult struct {
	ImageID    string
	Digest     string
	Size       int64
	PullTime   time.Duration
	Repository string
	Tag        string
}

// ImageTag represents an image tag in a repository
type ImageTag struct {
	Tag        string
	Digest     string
	Size       int64
	Created    time.Time
	Repository string
}

// RegistryConfig contains registry connection configuration
type RegistryConfig struct {
	Host     string
	Username string
	Password string
	Token    string
	Insecure bool
}
