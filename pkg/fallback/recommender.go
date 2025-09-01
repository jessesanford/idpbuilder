// Package fallback provides recommendation engine for certificate validation issues
package fallback

import (
	"fmt"
	"strings"
)

// Recommendation contains a suggested solution for a certificate problem
type Recommendation struct {
	Priority    int      // 1=highest priority, for ordering multiple recommendations
	Title       string   // Brief description of the solution
	Command     string   // Example command to execute (if applicable)
	Explanation string   // Detailed explanation of the solution
	Risks       []string // Security implications and risks
}

// Recommender generates actionable solutions for certificate validation problems
type Recommender interface {
	// Recommend generates prioritized solutions for a detected certificate problem
	Recommend(problem *CertProblem) ([]*Recommendation, error)
	
	// FormatRecommendations creates user-friendly output for display
	FormatRecommendations(recs []*Recommendation) string
	
	// GetQuickFix returns the highest priority recommendation for immediate action
	GetQuickFix(problem *CertProblem) (*Recommendation, error)
}

// DefaultRecommender implements Recommender with context-aware recommendations
type DefaultRecommender struct {
	registryURL     string // Target registry URL for context-specific advice
	insecureAllowed bool   // Whether --insecure flag is available/permitted
	environment     string // "development", "staging", "production"
}

// NewRecommender creates a recommender with registry context
func NewRecommender(registryURL string, allowInsecure bool) *DefaultRecommender {
	return &DefaultRecommender{
		registryURL:     registryURL,
		insecureAllowed: allowInsecure,
		environment:     "development", // Default assumption
	}
}

// NewRecommenderWithEnvironment creates a recommender with environment-specific advice
func NewRecommenderWithEnvironment(registryURL string, allowInsecure bool, environment string) *DefaultRecommender {
	return &DefaultRecommender{
		registryURL:     registryURL,
		insecureAllowed: allowInsecure,
		environment:     environment,
	}
}

// Recommend generates problem-specific recommendations with security considerations
func (r *DefaultRecommender) Recommend(problem *CertProblem) ([]*Recommendation, error) {
	if problem == nil {
		return nil, fmt.Errorf("problem cannot be nil")
	}

	var recommendations []*Recommendation

	switch problem.Type {
	case ProblemSelfSigned:
		recommendations = r.recommendForSelfSigned(problem)
	case ProblemExpired:
		recommendations = r.recommendForExpired(problem)
	case ProblemNotYetValid:
		recommendations = r.recommendForNotYetValid(problem)
	case ProblemHostnameMismatch:
		recommendations = r.recommendForHostnameMismatch(problem)
	case ProblemUntrustedCA:
		recommendations = r.recommendForUntrustedCA(problem)
	case ProblemUnknownAuthority:
		recommendations = r.recommendForUnknownAuthority(problem)
	default:
		recommendations = r.recommendForUnknown(problem)
	}

	return recommendations, nil
}

// recommendForSelfSigned provides solutions for self-signed certificate issues
func (r *DefaultRecommender) recommendForSelfSigned(problem *CertProblem) []*Recommendation {
	recs := make([]*Recommendation, 0)

	if r.environment == "development" || r.environment == "testing" {
		// Development environment - insecure flag is reasonable
		if r.insecureAllowed {
			recs = append(recs, &Recommendation{
				Priority:    1,
				Title:       "Use --insecure flag for development",
				Command:     "idpbuilder <command> --insecure",
				Explanation: "Skip TLS verification for self-signed certificates in development environments. This is the quickest solution for local Kind clusters.",
				Risks:       []string{"Not suitable for production", "Vulnerable to man-in-the-middle attacks"},
			})
		}
	}

	// Add certificate to trust store (always recommended for proper security)
	recs = append(recs, &Recommendation{
		Priority:    2,
		Title:       "Add certificate to system trust store",
		Command:     fmt.Sprintf("# Extract certificate:\necho -n | openssl s_client -connect %s:443 -servername %s 2>/dev/null | openssl x509 > registry.crt\n# Add to trust store (varies by OS)", r.extractHost(), r.extractHost()),
		Explanation: "Import the self-signed certificate into your system's trusted certificate store. This provides secure validation without disabling TLS.",
		Risks:       []string{"Requires system administrator privileges", "Must be done on each client system"},
	})

	// Registry-specific solution
	recs = append(recs, &Recommendation{
		Priority:    3,
		Title:       "Configure registry with proper TLS certificate",
		Command:     "# Generate proper certificate with valid CA or use Let's Encrypt",
		Explanation: "Replace the self-signed certificate with one issued by a trusted Certificate Authority. This is the most secure long-term solution.",
		Risks:       []string{"Requires access to registry configuration", "May require DNS setup for public CA"},
	})

	return recs
}

// recommendForExpired provides solutions for expired certificates
func (r *DefaultRecommender) recommendForExpired(problem *CertProblem) []*Recommendation {
	recs := make([]*Recommendation, 0)

	// Primary recommendation: renew certificate
	recs = append(recs, &Recommendation{
		Priority:    1,
		Title:       "Renew the expired certificate",
		Command:     "# Contact registry administrator to renew certificate",
		Explanation: "The certificate has expired and must be renewed. Contact your registry administrator or certificate authority.",
		Risks:       []string{"Service may be unavailable until renewal"},
	})

	// Temporary workaround for development
	if r.environment != "production" && r.insecureAllowed {
		recs = append(recs, &Recommendation{
			Priority:    2,
			Title:       "Temporary workaround with --insecure flag",
			Command:     "idpbuilder <command> --insecure",
			Explanation: "Temporarily bypass certificate validation while waiting for certificate renewal. Only use in development/testing.",
			Risks:       []string{"NEVER use in production", "Temporary solution only", "Security vulnerability"},
		})
	}

	return recs
}

// recommendForNotYetValid provides solutions for not-yet-valid certificates
func (r *DefaultRecommender) recommendForNotYetValid(problem *CertProblem) []*Recommendation {
	return []*Recommendation{
		{
			Priority:    1,
			Title:       "Check system clock synchronization",
			Command:     "# Linux: sudo ntpdate -s time.nist.gov\n# macOS: sudo sntp -sS time.apple.com",
			Explanation: "The certificate is not yet valid, which often indicates system clock skew. Synchronize your system time with NTP servers.",
			Risks:       []string{"May require system administrator privileges"},
		},
		{
			Priority:    2,
			Title:       "Verify certificate validity period",
			Command:     "echo -n | openssl s_client -connect " + r.extractHost() + ":443 -servername " + r.extractHost() + " 2>/dev/null | openssl x509 -noout -dates",
			Explanation: "Check the certificate's actual validity period to confirm if this is a time synchronization issue or a certificate configuration problem.",
			Risks:       []string{"None - informational only"},
		},
	}
}

// recommendForHostnameMismatch provides solutions for hostname validation failures
func (r *DefaultRecommender) recommendForHostnameMismatch(problem *CertProblem) []*Recommendation {
	recs := make([]*Recommendation, 0)

	// Extract valid hostnames from problem details
	var validHosts []string
	if hosts, ok := problem.Details["valid_hostnames"].([]string); ok {
		validHosts = hosts
	}

	if len(validHosts) > 0 {
		recs = append(recs, &Recommendation{
			Priority:    1,
			Title:       "Use correct hostname from certificate",
			Command:     fmt.Sprintf("# Use one of these hostnames: %s", strings.Join(validHosts, ", ")),
			Explanation: fmt.Sprintf("The certificate is valid for specific hostnames. Update your registry URL to use one of: %s", strings.Join(validHosts, ", ")),
			Risks:       []string{"May require DNS configuration or /etc/hosts changes"},
		})
	}

	// Add hostname mapping for development
	if r.environment == "development" && len(validHosts) > 0 {
		recs = append(recs, &Recommendation{
			Priority:    2,
			Title:       "Add hostname mapping for development",
			Command:     fmt.Sprintf("# Add to /etc/hosts:\n127.0.0.1 %s", validHosts[0]),
			Explanation: "Map the certificate's hostname to your local IP address for development testing.",
			Risks:       []string{"Requires root/administrator privileges", "Development only - not for production"},
		})
	}

	// Certificate update recommendation
	recs = append(recs, &Recommendation{
		Priority:    3,
		Title:       "Update certificate with correct Subject Alternative Names",
		Command:     "# Contact certificate authority to reissue with correct SANs",
		Explanation: "Request a new certificate that includes the hostname you're trying to use in the Subject Alternative Names.",
		Risks:       []string{"Requires certificate authority interaction", "May involve cost for new certificate"},
	})

	return recs
}

// recommendForUntrustedCA provides solutions for untrusted Certificate Authority issues
func (r *DefaultRecommender) recommendForUntrustedCA(problem *CertProblem) []*Recommendation {
	recs := make([]*Recommendation, 0)

	// Import CA certificate
	recs = append(recs, &Recommendation{
		Priority:    1,
		Title:       "Import root Certificate Authority certificate",
		Command:     "# Get CA certificate and add to system trust store",
		Explanation: "The certificate was issued by a CA that's not in your system's trust store. Import the root CA certificate.",
		Risks:       []string{"Requires system administrator privileges", "Ensure CA certificate is from trusted source"},
	})

	// Verify certificate chain
	recs = append(recs, &Recommendation{
		Priority:    2,
		Title:       "Verify complete certificate chain",
		Command:     "openssl s_client -connect " + r.extractHost() + ":443 -servername " + r.extractHost() + " -showcerts",
		Explanation: "Check if intermediate CA certificates are missing from the server configuration. The server should present the complete certificate chain.",
		Risks:       []string{"None - diagnostic command"},
	})

	// Insecure option for development
	if r.environment != "production" && r.insecureAllowed {
		recs = append(recs, &Recommendation{
			Priority:    3,
			Title:       "Development workaround with --insecure flag",
			Command:     "idpbuilder <command> --insecure",
			Explanation: "Skip CA validation for development/testing while resolving the trust store issue.",
			Risks:       []string{"NOT for production use", "Security vulnerability"},
		})
	}

	return recs
}

// recommendForUnknownAuthority provides solutions for unknown authority errors
func (r *DefaultRecommender) recommendForUnknownAuthority(problem *CertProblem) []*Recommendation {
	// Similar to untrusted CA but with focus on completely unknown CAs
	return r.recommendForUntrustedCA(problem)
}

// recommendForUnknown provides general solutions for unrecognized problems
func (r *DefaultRecommender) recommendForUnknown(problem *CertProblem) []*Recommendation {
	recs := make([]*Recommendation, 0)

	// General diagnostic
	recs = append(recs, &Recommendation{
		Priority:    1,
		Title:       "Manual certificate inspection",
		Command:     "echo -n | openssl s_client -connect " + r.extractHost() + ":443 -servername " + r.extractHost() + " 2>/dev/null | openssl x509 -text -noout",
		Explanation: "Manually inspect the certificate to understand the validation failure.",
		Risks:       []string{"None - diagnostic command"},
	})

	// Debug logging
	recs = append(recs, &Recommendation{
		Priority:    2,
		Title:       "Enable debug logging",
		Command:     "# Add debug flags to see detailed error information",
		Explanation: "Enable verbose logging to get more details about the certificate validation failure.",
		Risks:       []string{"None - informational only"},
	})

	return recs
}

// GetQuickFix returns the highest priority recommendation for immediate action
func (r *DefaultRecommender) GetQuickFix(problem *CertProblem) (*Recommendation, error) {
	recommendations, err := r.Recommend(problem)
	if err != nil {
		return nil, err
	}
	
	if len(recommendations) == 0 {
		return nil, fmt.Errorf("no recommendations available for problem type: %s", problem.Type)
	}
	
	// Find highest priority (lowest number)
	quickFix := recommendations[0]
	for _, rec := range recommendations {
		if rec.Priority < quickFix.Priority {
			quickFix = rec
		}
	}
	
	return quickFix, nil
}

// FormatRecommendations creates user-friendly output for displaying recommendations
func (r *DefaultRecommender) FormatRecommendations(recs []*Recommendation) string {
	if len(recs) == 0 {
		return "No recommendations available."
	}

	var sb strings.Builder
	
	sb.WriteString("Certificate Problem Solutions\n")
	sb.WriteString("============================\n\n")

	for i, rec := range recs {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, rec.Title))
		sb.WriteString("   " + strings.ReplaceAll(rec.Explanation, "\n", "\n   ") + "\n\n")
		
		if rec.Command != "" {
			sb.WriteString("   Command:\n")
			commandLines := strings.Split(rec.Command, "\n")
			for _, line := range commandLines {
				if strings.TrimSpace(line) != "" {
					sb.WriteString("   $ " + line + "\n")
				}
			}
			sb.WriteString("\n")
		}
		
		if len(rec.Risks) > 0 {
			sb.WriteString("   Security Considerations:\n")
			for _, risk := range rec.Risks {
				sb.WriteString("   ⚠ " + risk + "\n")
			}
			sb.WriteString("\n")
		}
		
		if i < len(recs)-1 {
			sb.WriteString("---\n\n")
		}
	}

	return sb.String()
}

// extractHost extracts hostname from registry URL for command examples
func (r *DefaultRecommender) extractHost() string {
	if r.registryURL == "" {
		return "registry.example.com"
	}
	
	// Simple extraction - remove protocol and path
	url := r.registryURL
	if strings.HasPrefix(url, "https://") {
		url = url[8:]
	} else if strings.HasPrefix(url, "http://") {
		url = url[7:]
	}
	
	// Remove path
	if idx := strings.Index(url, "/"); idx > 0 {
		url = url[:idx]
	}
	
	// Remove port for hostname examples
	if idx := strings.Index(url, ":"); idx > 0 {
		url = url[:idx]
	}
	
	return url
}