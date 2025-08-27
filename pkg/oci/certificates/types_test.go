package certificates

import "testing"

func TestCertificateStatus_Constants(t *testing.T) {
	if string(CertificateStatusActive) != "active" {
		t.Error("CertificateStatusActive mismatch")
	}
	if string(CertificateStatusExpired) != "expired" {
		t.Error("CertificateStatusExpired mismatch")
	}
	if string(CertificateStatusPending) != "pending" {
		t.Error("CertificateStatusPending mismatch")
	}
}

func TestEventType_Constants(t *testing.T) {
	if string(EventCertificateAdded) != "certificate_added" {
		t.Error("EventCertificateAdded mismatch")
	}
	if string(EventValidationFailed) != "validation_failed" {
		t.Error("EventValidationFailed mismatch")
	}
	if string(EventPoolUpdated) != "pool_updated" {
		t.Error("EventPoolUpdated mismatch")
	}
}

func TestCertificate_Creation(t *testing.T) {
	cert := &Certificate{
		ID:     "test-cert-1",
		Status: CertificateStatusActive,
		Tags:   map[string]string{"env": "test"},
	}

	if cert.ID != "test-cert-1" || cert.Status != CertificateStatusActive || cert.Tags["env"] != "test" {
		t.Error("Certificate field mismatch")
	}
}

func TestEvent_Creation(t *testing.T) {
	event := &Event{
		Type:          EventCertificateAdded,
		CertificateID: "test-cert-1",
		Metadata:      map[string]interface{}{"source": "api"},
	}

	if event.Type != EventCertificateAdded || event.CertificateID != "test-cert-1" || event.Metadata["source"] != "api" {
		t.Error("Event field mismatch")
	}
}

func TestValidationRule_Creation(t *testing.T) {
	rule := &ValidationRule{Name: "expiration-check", Enabled: true, Critical: true}
	if rule.Name != "expiration-check" || !rule.Enabled || !rule.Critical {
		t.Error("ValidationRule field mismatch")
	}
}

func TestValidationResult_Creation(t *testing.T) {
	result := &ValidationResult{
		Valid:  false,
		Errors: []ValidationError{NewValidationError("CERT_EXPIRED", "Certificate expired")},
	}
	if result.Valid || len(result.Errors) != 1 || result.Errors[0].Code != "CERT_EXPIRED" {
		t.Error("ValidationResult field mismatch")
	}
}