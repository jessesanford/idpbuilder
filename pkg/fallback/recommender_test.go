package fallback

import (
	"testing"
)

func TestNewRecommender(t *testing.T) {
	recommender := NewRecommender("https://registry.example.com", true)

	if recommender == nil {
		t.Fatal("NewRecommender returned nil")
	}

	if recommender.registryURL != "https://registry.example.com" {
		t.Error("Registry URL not set correctly")
	}

	if !recommender.insecureAllowed {
		t.Error("Insecure allowed not set correctly")
	}

	if recommender.environment != "development" {
		t.Error("Default environment not set to development")
	}
}

func TestRecommend_NilProblem(t *testing.T) {
	recommender := NewRecommender("https://registry.example.com", true)

	recs, err := recommender.Recommend(nil)

	if err == nil {
		t.Error("Expected error for nil problem")
	}

	if recs != nil {
		t.Error("Expected nil recommendations for nil problem")
	}
}

func TestRecommend_SelfSigned(t *testing.T) {
	recommender := NewRecommender("https://registry.example.com", true)
	problem := &CertProblem{
		Type: ProblemSelfSigned,
	}

	recs, err := recommender.Recommend(problem)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(recs) == 0 {
		t.Error("Expected recommendations for self-signed problem")
	}

	// Check that the first recommendation is about using --insecure
	if len(recs) > 0 && recs[0].Priority != 1 {
		t.Error("Expected first recommendation to have priority 1")
	}
}

func TestRecommend_Expired(t *testing.T) {
	recommender := NewRecommender("https://registry.example.com", true)
	problem := &CertProblem{
		Type: ProblemExpired,
	}

	recs, err := recommender.Recommend(problem)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(recs) == 0 {
		t.Error("Expected recommendations for expired problem")
	}
}

func TestGetQuickFix(t *testing.T) {
	recommender := NewRecommender("https://registry.example.com", true)
	problem := &CertProblem{
		Type: ProblemSelfSigned,
	}

	quickFix, err := recommender.GetQuickFix(problem)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if quickFix == nil {
		t.Fatal("Expected quick fix recommendation")
	}

	if quickFix.Priority != 1 {
		t.Error("Expected quick fix to have priority 1")
	}
}

func TestFormatRecommendations_Empty(t *testing.T) {
	recommender := NewRecommender("https://registry.example.com", true)

	formatted := recommender.FormatRecommendations([]*Recommendation{})

	if formatted != "No recommendations available." {
		t.Errorf("Expected empty message, got: %s", formatted)
	}
}

func TestFormatRecommendations_WithData(t *testing.T) {
	recommender := NewRecommender("https://registry.example.com", true)
	recs := []*Recommendation{
		{
			Priority:    1,
			Title:       "Test Solution",
			Command:     "test command",
			Explanation: "Test explanation",
			Risks:       []string{"Test risk"},
		},
	}

	formatted := recommender.FormatRecommendations(recs)

	if formatted == "" {
		t.Error("Expected formatted output")
	}

	// Check that the formatted output contains key elements
	if !contains(formatted, "Test Solution") {
		t.Error("Expected formatted output to contain title")
	}

	if !contains(formatted, "Test explanation") {
		t.Error("Expected formatted output to contain explanation")
	}
}

func TestExtractHost(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://registry.example.com", "registry.example.com"},
		{"http://registry.example.com", "registry.example.com"},
		{"registry.example.com:5000", "registry.example.com"},
		{"https://registry.example.com:5000/path", "registry.example.com"},
		{"", "registry.example.com"}, // default
	}

	for _, test := range tests {
		recommender := NewRecommender(test.input, true)
		result := recommender.extractHost()

		if result != test.expected {
			t.Errorf("For input %s, expected: %s, got: %s", test.input, test.expected, result)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
