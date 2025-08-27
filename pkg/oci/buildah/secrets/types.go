package secrets

import (
	"context"
	"regexp"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

// SecretManager handles secure build args and secrets for Buildah operations
type SecretManager struct {
	k8sClient     kubernetes.Interface
	namespace     string
	tempDir       string
	logger        *logrus.Logger
	secretFiles   map[string]*secretFile
	mutex         sync.RWMutex
	cleanupTicker *time.Ticker
	ctx           context.Context
	cancel        context.CancelFunc
}

// secretFile represents a temporary secret file
type secretFile struct {
	Path      string
	CreatedAt time.Time
	TTL       time.Duration
}

// BuildArgs represents sanitized build arguments
type BuildArgs struct {
	Args      map[string]string
	SecretIDs []string
}

// SecretMount represents a secret to be mounted during build
type SecretMount struct {
	ID      string
	Source  string
	Target  string
	Type    SecretType
	Options map[string]string
}

// SecretType defines the type of secret source
type SecretType string

const (
	SecretTypeFile       SecretType = "file"
	SecretTypeKubernetes SecretType = "kubernetes"
	SecretTypeVault      SecretType = "vault"
	SecretTypeEnv        SecretType = "env"
)

// Patterns to redact in logs
var secretPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(password|passwd|secret|token|key|credential)[\s]*[=:][\s]*["\']?([^\s"']+)`),
	regexp.MustCompile(`(?i)(bearer\s+)([a-zA-Z0-9\-._~+/]+=*)`),
	regexp.MustCompile(`(?i)(api[_-]?key|access[_-]?token)[\s]*[=:][\s]*["\']?([^\s"']+)`),
}