package types

// Secret represents secret data for printing/display purposes
type Secret struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	IsCore    bool              `json:"isCore"`
	Username  string            `json:"username,omitempty"`
	Password  string            `json:"password,omitempty"`
	Token     string            `json:"token,omitempty"`
	Data      map[string]string `json:"data,omitempty"`
}
