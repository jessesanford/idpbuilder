// Copyright 2024 idpbuilder Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package registry

import (
	"crypto/tls"
	"net/http"
)

// TLSConfig creates a TLS configuration based on registry configuration.
// This handles the common case of local development with self-signed certificates.
func TLSConfig(config RegistryConfig) *tls.Config {
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Handle insecure connections for development
	if config.Insecure || config.SkipTLSVerify {
		tlsConfig.InsecureSkipVerify = true
	}

	return tlsConfig
}

// HTTPTransport creates an HTTP transport with appropriate TLS configuration.
// This is used by the registry client for HTTP operations.
func HTTPTransport(config RegistryConfig) *http.Transport {
	return &http.Transport{
		TLSClientConfig: TLSConfig(config),
	}
}

// HTTPClient creates an HTTP client configured for registry operations.
func HTTPClient(config RegistryConfig) *http.Client {
	return &http.Client{
		Transport: HTTPTransport(config),
		Timeout:   config.Timeout,
	}
}