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

// Recommend generates problem-specific recommendations with security considerations
func (r *DefaultRecommender) Recommend(problem *CertProblem) ([]*Recommendation, error) {
	if problem == nil {
		return nil, fmt.Errorf("problem cannot be nil")
	}

	switch problem.Type {
	case ProblemSelfSigned:
		return r.getRecommendations([]string{"Use --insecure for dev", "Add cert to trust store"},
			[]string{"idpbuilder <command> --insecure", "openssl s_client -connect " + r.extractHost() + ":443"})
	case ProblemExpired:
		return r.getRecommendations([]string{"Renew certificate", "Use --insecure temporarily"},
			[]string{"# Contact administrator", "idpbuilder <command> --insecure"})
	case ProblemNotYetValid:
		return r.getRecommendations([]string{"Check system clock", "Verify cert dates"},
			[]string{"sudo ntpdate -s time.nist.gov", "openssl x509 -noout -dates"})
	case ProblemHostnameMismatch:
		validHosts := ""
		if hosts, ok := problem.Details["valid_hostnames"].([]string); ok && len(hosts) > 0 {
			validHosts = strings.Join(hosts, ", ")
		}
		return r.getRecommendations([]string{"Use correct hostname", "Update certificate"},
			[]string{"# Use: " + validHosts, "# Request new cert with correct SANs"})
	case ProblemUntrustedCA:
		return r.getRecommendations([]string{"Import CA certificate", "Verify cert chain"},
			[]string{"# Add CA to trust store", "openssl s_client -connect " + r.extractHost() + ":443 -showcerts"})
	default:
		return r.getRecommendations([]string{"Inspect certificate manually", "Enable debug logging"},
			[]string{"openssl x509 -text -noout", "# Add debug flags"})
	}
}

// getRecommendations creates recommendations from titles and commands
func (r *DefaultRecommender) getRecommendations(titles []string, commands []string) ([]*Recommendation, error) {
	recs := make([]*Recommendation, 0, len(titles))
	for i, title := range titles {
		cmd := ""
		if i < len(commands) {
			cmd = commands[i]
		}

		// Add --insecure option for development if applicable
		risks := []string{"Check documentation for full details"}
		if strings.Contains(cmd, "--insecure") && r.environment != "development" {
			risks = []string{"NOT for production use", "Security vulnerability"}
		}

		recs = append(recs, &Recommendation{
			Priority:    i + 1,
			Title:       title,
			Command:     cmd,
			Explanation: "Recommended solution for this certificate problem.",
			Risks:       risks,
		})
	}
	return recs, nil
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
	return recommendations[0], nil // First is highest priority
}

// FormatRecommendations creates user-friendly output for displaying recommendations
func (r *DefaultRecommender) FormatRecommendations(recs []*Recommendation) string {
	if len(recs) == 0 {
		return "No recommendations available."
	}

	var sb strings.Builder
	sb.WriteString("Certificate Problem Solutions\n============================\n\n")

	for i, rec := range recs {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, rec.Title))
		sb.WriteString("   " + rec.Explanation + "\n")

		if rec.Command != "" {
			sb.WriteString("   Command: " + rec.Command + "\n")
		}

		if len(rec.Risks) > 0 {
			sb.WriteString("   Risks: " + strings.Join(rec.Risks, ", ") + "\n")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// extractHost extracts hostname from registry URL for command examples
func (r *DefaultRecommender) extractHost() string {
	if r.registryURL == "" {
		return "registry.example.com"
	}

	url := r.registryURL
	if strings.HasPrefix(url, "https://") {
		url = url[8:]
	} else if strings.HasPrefix(url, "http://") {
		url = url[7:]
	}

	if idx := strings.Index(url, "/"); idx > 0 {
		url = url[:idx]
	}
	if idx := strings.Index(url, ":"); idx > 0 {
		url = url[:idx]
	}
	return url
}
