package certs

import (
	"time"
)

// ChainValidationResult contains detailed chain validation results
type ChainValidationResult struct {
	Valid            bool                  `json:"valid"`
	ChainComplete    bool                  `json:"chain_complete"`
	TrustAnchorFound bool                  `json:"trust_anchor_found"`
	ChainLength      int                   `json:"chain_length"`
	Certificates     []CertificateSummary  `json:"certificates"`
	Issues           []ValidationIssue     `json:"issues"`
}

// CertificateSummary provides a summary of a certificate in the chain
type CertificateSummary struct {
	Subject      string    `json:"subject"`
	Issuer       string    `json:"issuer"`
	SerialNumber string    `json:"serial_number"`
	NotBefore    time.Time `json:"not_before"`
	NotAfter     time.Time `json:"not_after"`
	IsCA         bool      `json:"is_ca"`
	Position     int       `json:"position"` // Position in chain (0 = leaf)
}

// ValidationIssue represents a specific validation problem
type ValidationIssue struct {
	Severity    IssueSeverity `json:"severity"`
	Code        string        `json:"code"`
	Message     string        `json:"message"`
	Certificate string        `json:"certificate"` // Which cert in chain has the issue
	Remediation string        `json:"remediation"` // Suggested fix
}

// IssueSeverity defines the severity of validation issues
type IssueSeverity int

const (
	SeverityWarning IssueSeverity = iota
	SeverityError
	SeverityCritical
)

// String returns string representation of severity
func (s IssueSeverity) String() string {
	switch s {
	case SeverityWarning:
		return "WARNING"
	case SeverityError:
		return "ERROR"
	case SeverityCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

// ChainExpiryResult contains expiry information for the entire chain
type ChainExpiryResult struct {
	ChainValid      bool                  `json:"chain_valid"`
	ExpiringCerts   []ExpiringCertificate `json:"expiring_certs"`
	ExpiredCerts    []ExpiredCertificate  `json:"expired_certs"`
	MinDaysToExpiry int                   `json:"min_days_to_expiry"`
}

// ExpiringCertificate represents a certificate that will expire soon
type ExpiringCertificate struct {
	Subject         string    `json:"subject"`
	Position        int       `json:"position"`
	DaysUntilExpiry int       `json:"days_until_expiry"`
	ExpiryDate      time.Time `json:"expiry_date"`
}

// ExpiredCertificate represents a certificate that has already expired
type ExpiredCertificate struct {
	Subject      string    `json:"subject"`
	Position     int       `json:"position"`
	DaysExpired  int       `json:"days_expired"`
	ExpiredDate  time.Time `json:"expired_date"`
}

// CertDiagnosticsReport provides comprehensive diagnostic information
type CertDiagnosticsReport struct {
	Timestamp          time.Time             `json:"timestamp"`
	Hostname           string                `json:"hostname"`
	CertificateDetails *CertificateDetails   `json:"certificate_details"`
	ChainAnalysis      *ChainAnalysis        `json:"chain_analysis"`
	HostnameValidation *HostnameValidation   `json:"hostname_validation"`
	TrustStoreAnalysis *TrustStoreAnalysis   `json:"trust_store_analysis"`
	Recommendations    []Recommendation      `json:"recommendations"`
}

// CertificateDetails contains detailed information about the main certificate
type CertificateDetails struct {
	Subject           string               `json:"subject"`
	Issuer            string               `json:"issuer"`
	SerialNumber      string               `json:"serial_number"`
	NotBefore         time.Time            `json:"not_before"`
	NotAfter          time.Time            `json:"not_after"`
	IsCA              bool                 `json:"is_ca"`
	IsSelfSigned      bool                 `json:"is_self_signed"`
	DNSNames          []string             `json:"dns_names"`
	IPAddresses       []string             `json:"ip_addresses"`
	KeyUsage          []string             `json:"key_usage"`
	ExtKeyUsage       []string             `json:"ext_key_usage"`
	SignatureAlgorithm string              `json:"signature_algorithm"`
	PublicKeyAlgorithm string              `json:"public_key_algorithm"`
	ValidationResult  *ValidationResult    `json:"validation_result"`
}

// ChainAnalysis contains information about the certificate chain
type ChainAnalysis struct {
	ChainLength       int                  `json:"chain_length"`
	ChainValid        bool                 `json:"chain_valid"`
	TrustAnchorFound  bool                 `json:"trust_anchor_found"`
	Certificates      []CertificateSummary `json:"certificates"`
	ChainIssues       []ValidationIssue    `json:"chain_issues"`
}

// HostnameValidation contains hostname verification results
type HostnameValidation struct {
	Hostname     string `json:"hostname"`
	Valid        bool   `json:"valid"`
	MatchType    string `json:"match_type"` // "exact", "wildcard", "san", or "none"
	MatchedValue string `json:"matched_value"`
	Error        string `json:"error,omitempty"`
}

// TrustStoreAnalysis contains trust store verification results
type TrustStoreAnalysis struct {
	TrustStoreChecked bool   `json:"trust_store_checked"`
	CertificateTrusted bool  `json:"certificate_trusted"`
	TrustPath         []string `json:"trust_path"`
	TrustError        string   `json:"trust_error,omitempty"`
}

// Recommendation provides actionable advice for resolving issues
type Recommendation struct {
	Priority    RecommendationPriority `json:"priority"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Command     string                 `json:"command"`     // Specific command to run
	Link        string                 `json:"link"`        // Documentation link
}

// RecommendationPriority defines the priority of recommendations
type RecommendationPriority int

const (
	PriorityLow RecommendationPriority = iota
	PriorityMedium
	PriorityHigh
	PriorityCritical
)

// String returns string representation of priority
func (p RecommendationPriority) String() string {
	switch p {
	case PriorityLow:
		return "LOW"
	case PriorityMedium:
		return "MEDIUM"
	case PriorityHigh:
		return "HIGH"
	case PriorityCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}