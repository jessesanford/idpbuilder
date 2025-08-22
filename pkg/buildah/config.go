package buildah

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/cnoe-io/idpbuilder/api/v1alpha1"
)

// ApplyConfiguration applies build customization settings to build options
func (c *client) ApplyConfiguration(opts *BuildOptions, config v1alpha1.BuildCustomizationSpec) error {
	c.logger.V(1).Info("Applying build configuration", 
		"protocol", config.Protocol,
		"host", config.Host,
		"port", config.Port)

	// Apply registry configuration if specified
	if err := c.applyRegistryConfig(opts, config); err != nil {
		return fmt.Errorf("failed to apply registry configuration: %w", err)
	}

	// Apply network configuration
	if err := c.applyNetworkConfig(opts, config); err != nil {
		return fmt.Errorf("failed to apply network configuration: %w", err)
	}

	// Apply build-specific configuration
	if err := c.applyBuildConfig(opts, config); err != nil {
		return fmt.Errorf("failed to apply build configuration: %w", err)
	}

	return nil
}

// applyRegistryConfig configures registry-related build settings
func (c *client) applyRegistryConfig(opts *BuildOptions, config v1alpha1.BuildCustomizationSpec) error {
	// If host is specified, it might be used as a registry
	if config.Host != "" {
		// Construct registry URL from host and port
		registryURL := config.Host
		if config.Port != "" {
			registryURL = registryURL + ":" + config.Port
		}

		// Add protocol prefix if specified
		if config.Protocol != "" {
			registryURL = config.Protocol + "://" + registryURL
		}

		c.logger.V(1).Info("Configured registry", "url", registryURL)

		// Set as default registry for the client
		c.defaultRegistry = registryURL

		// If the image tag doesn't already include a registry, add it
		if opts.Tag != "" && !strings.Contains(opts.Tag, "/") {
			opts.Tag = registryURL + "/" + opts.Tag
			c.logger.V(1).Info("Updated image tag with registry", "tag", opts.Tag)
		}
	}

	return nil
}

// applyNetworkConfig configures network-related build settings
func (c *client) applyNetworkConfig(opts *BuildOptions, config v1alpha1.BuildCustomizationSpec) error {
	if opts.BuildArgs == nil {
		opts.BuildArgs = make(map[string]string)
	}

	// Configure HTTP/HTTPS proxy settings based on protocol
	if config.Protocol != "" {
		switch strings.ToLower(config.Protocol) {
		case "http":
			opts.BuildArgs["HTTP_PROXY"] = buildProxyURL(config, false)
		case "https":
			opts.BuildArgs["HTTPS_PROXY"] = buildProxyURL(config, true)
			opts.BuildArgs["HTTP_PROXY"] = buildProxyURL(config, false)
		}
	}

	// Configure ingress host if specified
	if config.IngressHost != "" {
		opts.BuildArgs["INGRESS_HOST"] = config.IngressHost
		c.logger.V(1).Info("Set ingress host build arg", "host", config.IngressHost)
	}

	return nil
}

// applyBuildConfig applies general build configuration settings
func (c *client) applyBuildConfig(opts *BuildOptions, config v1alpha1.BuildCustomizationSpec) error {
	if opts.Labels == nil {
		opts.Labels = make(map[string]string)
	}

	// Add configuration labels to the built image
	opts.Labels["idpbuilder.cnoe.io/protocol"] = config.Protocol
	opts.Labels["idpbuilder.cnoe.io/host"] = config.Host
	
	if config.IngressHost != "" {
		opts.Labels["idpbuilder.cnoe.io/ingress-host"] = config.IngressHost
	}
	
	if config.Port != "" {
		opts.Labels["idpbuilder.cnoe.io/port"] = config.Port
	}

	// Configure path routing if enabled
	if config.UsePathRouting {
		opts.BuildArgs["USE_PATH_ROUTING"] = "true"
		opts.Labels["idpbuilder.cnoe.io/path-routing"] = "enabled"
		c.logger.V(1).Info("Enabled path routing configuration")
	}

	// Configure static password mode if enabled
	if config.StaticPassword {
		opts.BuildArgs["STATIC_PASSWORD"] = "true"
		opts.Labels["idpbuilder.cnoe.io/static-password"] = "enabled"
		c.logger.V(1).Info("Enabled static password configuration")
	}

	// Add self-signed certificate information if available
	if config.SelfSignedCert != "" {
		// Don't include the actual certificate in build args for security
		opts.Labels["idpbuilder.cnoe.io/tls"] = "self-signed"
		c.logger.V(1).Info("TLS configuration detected")
	}

	c.logger.V(1).Info("Applied build configuration", 
		"build_args_count", len(opts.BuildArgs),
		"labels_count", len(opts.Labels))

	return nil
}

// buildProxyURL constructs a proxy URL from build customization config
func buildProxyURL(config v1alpha1.BuildCustomizationSpec, https bool) string {
	if config.Host == "" {
		return ""
	}

	scheme := "http"
	if https && config.Protocol == "https" {
		scheme = "https"
	}

	proxyURL := &url.URL{
		Scheme: scheme,
		Host:   config.Host,
	}

	if config.Port != "" {
		proxyURL.Host = config.Host + ":" + config.Port
	}

	return proxyURL.String()
}

// ValidateConfiguration validates that the build customization is compatible with buildah
func ValidateConfiguration(config v1alpha1.BuildCustomizationSpec) error {
	// Validate protocol
	if config.Protocol != "" {
		validProtocols := []string{"http", "https"}
		if !contains(validProtocols, strings.ToLower(config.Protocol)) {
			return fmt.Errorf("unsupported protocol: %s (supported: %s)", 
				config.Protocol, strings.Join(validProtocols, ", "))
		}
	}

	// Validate host format
	if config.Host != "" {
		if strings.Contains(config.Host, "://") {
			return fmt.Errorf("host should not include protocol scheme: %s", config.Host)
		}
	}

	// Validate port format
	if config.Port != "" {
		if !isValidPort(config.Port) {
			return fmt.Errorf("invalid port format: %s", config.Port)
		}
	}

	// Validate ingress host format  
	if config.IngressHost != "" {
		if strings.Contains(config.IngressHost, "://") {
			return fmt.Errorf("ingress host should not include protocol scheme: %s", config.IngressHost)
		}
	}

	return nil
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// isValidPort checks if a string represents a valid port number
func isValidPort(port string) bool {
	if port == "" {
		return false
	}

	// Simple validation - just check if it's numeric and in valid range
	for _, r := range port {
		if r < '0' || r > '9' {
			return false
		}
	}

	// Convert to int and check range (1-65535)
	if len(port) > 5 {
		return false
	}

	var portNum int
	for _, r := range port {
		portNum = portNum*10 + int(r-'0')
	}

	return portNum >= 1 && portNum <= 65535
}