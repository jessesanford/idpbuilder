package registry

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/oci/api"
	"github.com/cnoe-io/idpbuilder/pkg/oci/security"
)

type registryClient struct {
	httpClient    *http.Client
	securityMgr   api.SecurityManager
	defaultAuth   api.AuthConfig
	transportOpts TransportOptions
}

type TransportOptions struct {
	MaxRetries         int
	RetryBackoff       time.Duration
	InsecureSkipVerify bool
	CACertPath         string
	ClientCertPath     string
	ClientKeyPath      string
	Timeout            time.Duration
}

type Option func(*registryClient)

func WithSecurityManager(sm api.SecurityManager) Option {
	return func(rc *registryClient) { rc.securityMgr = sm }
}

func WithDefaultAuth(auth api.AuthConfig) Option {
	return func(rc *registryClient) { rc.defaultAuth = auth }
}

func WithTransportOptions(opts TransportOptions) Option {
	return func(rc *registryClient) { rc.transportOpts = opts }
}

func WithInsecureSkipVerify(insecure bool) Option {
	return func(rc *registryClient) { rc.transportOpts.InsecureSkipVerify = insecure }
}

func WithCACertificate(caCertPath string) Option {
	return func(rc *registryClient) { rc.transportOpts.CACertPath = caCertPath }
}

func NewRegistryClient(opts ...Option) api.RegistryClient {
	rc := &registryClient{
		transportOpts: TransportOptions{
			MaxRetries:   3,
			RetryBackoff: 2 * time.Second,
			Timeout:      30 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(rc)
	}
	rc.httpClient = rc.createHTTPClient()
	return rc
}

func (rc *registryClient) createHTTPClient() *http.Client {
	tlsConfig := &tls.Config{InsecureSkipVerify: rc.transportOpts.InsecureSkipVerify}

	if rc.transportOpts.CACertPath != "" {
		if caCert, err := ioutil.ReadFile(rc.transportOpts.CACertPath); err == nil {
			caCertPool := x509.NewCertPool()
			if caCertPool.AppendCertsFromPEM(caCert) {
				tlsConfig.RootCAs = caCertPool
				tlsConfig.InsecureSkipVerify = false
			}
		}
	}

	if rc.transportOpts.ClientCertPath != "" && rc.transportOpts.ClientKeyPath != "" {
		if cert, err := tls.LoadX509KeyPair(rc.transportOpts.ClientCertPath, rc.transportOpts.ClientKeyPath); err == nil {
			tlsConfig.Certificates = []tls.Certificate{cert}
		}
	}

	baseTransport := &http.Transport{
		TLSClientConfig:       tlsConfig,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &http.Client{
		Transport: newRetryTransport(baseTransport, rc.transportOpts.MaxRetries, rc.transportOpts.RetryBackoff),
		Timeout:   rc.transportOpts.Timeout,
	}
}

func (rc *registryClient) isGiteaRegistry(serverAddress string) bool {
	return strings.Contains(serverAddress, "gitea.cnoe.localtest.me")
}

func (rc *registryClient) getRegistryURL(ref *imageReference) string {
	if ref.Registry == "" {
		return "https://registry-1.docker.io"
	}
	if rc.isGiteaRegistry(ref.Registry) {
		return fmt.Sprintf("https://%s", ref.Registry)
	}
	if !strings.HasPrefix(ref.Registry, "http://") && !strings.HasPrefix(ref.Registry, "https://") {
		return fmt.Sprintf("https://%s", ref.Registry)
	}
	return ref.Registry
}

func (rc *registryClient) getRegistryBaseURL() string {
	if rc.defaultAuth.ServerAddress == "" {
		return "https://registry-1.docker.io"
	}
	return rc.getRegistryURL(&imageReference{Registry: rc.defaultAuth.ServerAddress})
}

func (rc *registryClient) Close() error {
	if rc.httpClient != nil {
		rc.httpClient.CloseIdleConnections()
	}
	return nil
}

type imageReference struct {
	Registry   string
	Namespace  string
	Repository string
	Name       string
	Tag        string
	Digest     string
}

func parseImageReference(image string) (*imageReference, error) {
	ref := &imageReference{}
	if strings.Contains(image, "@") {
		parts := strings.SplitN(image, "@", 2)
		image = parts[0]
		ref.Digest = parts[1]
	}
	if strings.Contains(image, ":") && ref.Digest == "" {
		parts := strings.SplitN(image, ":", 2)
		image = parts[0]
		ref.Tag = parts[1]
	} else if ref.Digest == "" {
		ref.Tag = "latest"
	}
	parts := strings.Split(image, "/")
	switch len(parts) {
	case 1:
		ref.Repository = parts[0]
		ref.Name = parts[0]
		ref.Registry = "registry-1.docker.io"
		ref.Namespace = "library"
	case 2:
		if strings.Contains(parts[0], ".") || strings.Contains(parts[0], ":") {
			ref.Registry = parts[0]
			ref.Repository = parts[1]
			ref.Name = parts[1]
		} else {
			ref.Registry = "registry-1.docker.io"
			ref.Namespace = parts[0]
			ref.Repository = parts[1]
			ref.Name = parts[1]
		}
	default:
		ref.Registry = parts[0]
		ref.Namespace = parts[1]
		ref.Repository = strings.Join(parts[2:], "/")
		ref.Name = strings.Join(parts[1:], "/")
	}
	return ref, nil
}

func (ref *imageReference) String() string {
	var result strings.Builder
	if ref.Registry != "" && ref.Registry != "registry-1.docker.io" {
		result.WriteString(ref.Registry + "/")
	}
	if ref.Namespace != "" && ref.Namespace != "library" {
		result.WriteString(ref.Namespace + "/")
	}
	result.WriteString(ref.Repository)
	if ref.Tag != "" && ref.Digest == "" {
		result.WriteString(":" + ref.Tag)
	}
	if ref.Digest != "" {
		result.WriteString("@" + ref.Digest)
	}
	return result.String()
}