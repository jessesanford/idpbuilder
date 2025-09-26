package buildah

import (
	"fmt"
	"testing"
)

// MockBuildContextManager implements BuildContextManager for testing
type MockBuildContextManager struct {
	contexts  map[string]bool
	artifacts map[string]map[string]string
	current   string
}

func NewMockBuildContextManager() *MockBuildContextManager {
	return &MockBuildContextManager{
		contexts:  make(map[string]bool),
		artifacts: make(map[string]map[string]string),
	}
}

func (m *MockBuildContextManager) CreateStageContext(stageName string) error {
	m.contexts[stageName] = true
	if m.artifacts[stageName] == nil {
		m.artifacts[stageName] = make(map[string]string)
	}
	return nil
}

func (m *MockBuildContextManager) SetCurrentContext(stageName string) error {
	if !m.contexts[stageName] {
		return fmt.Errorf("context %s does not exist", stageName)
	}
	m.current = stageName
	return nil
}

func (m *MockBuildContextManager) GetArtifacts(stageName string) (map[string]string, error) {
	if artifacts, exists := m.artifacts[stageName]; exists {
		return artifacts, nil
	}
	return make(map[string]string), nil
}

func (m *MockBuildContextManager) PreserveArtifact(stageName, path, alias string) error {
	if m.artifacts[stageName] == nil {
		m.artifacts[stageName] = make(map[string]string)
	}
	m.artifacts[stageName][path] = alias
	return nil
}

func TestMultiStageBuilder_ParseDockerfile(t *testing.T) {
	tests := []struct {
		name       string
		dockerfile string
		wantStages int
		wantError  bool
		checkStage func(*MultiStageBuilder) bool
	}{
		{
			name: "simple multi-stage",
			dockerfile: `FROM golang:1.20 AS builder
RUN go build -o app
FROM alpine:3.18
COPY --from=builder /app /app`,
			wantStages: 2,
			wantError:  false,
			checkStage: func(m *MultiStageBuilder) bool {
				stage, err := m.GetStage("builder")
				if err != nil {
					return false
				}
				return stage.BaseImage == "golang:1.20" && len(stage.Instructions) == 1
			},
		},
		{
			name: "named stages with dependencies",
			dockerfile: `FROM node:18 AS frontend
RUN npm build
FROM golang:1.20 AS backend
RUN go build
FROM alpine:3.18 AS final
COPY --from=frontend /dist /static
COPY --from=backend /app /app`,
			wantStages: 3,
			wantError:  false,
			checkStage: func(m *MultiStageBuilder) bool {
				stage, err := m.GetStage("final")
				if err != nil {
					return false
				}
				return len(stage.Dependencies) == 2 &&
					contains(stage.Dependencies, "frontend") &&
					contains(stage.Dependencies, "backend")
			},
		},
		{
			name: "auto-generated stage names",
			dockerfile: `FROM golang:1.20
RUN go build
FROM alpine:3.18
COPY --from=stage-0 /app /app`,
			wantStages: 2,
			wantError:  false,
			checkStage: func(m *MultiStageBuilder) bool {
				_, err1 := m.GetStage("stage-0")
				_, err2 := m.GetStage("stage-1")
				return err1 == nil && err2 == nil
			},
		},
		{
			name:       "empty dockerfile",
			dockerfile: "",
			wantStages: 0,
			wantError:  true,
		},
		{
			name: "dockerfile with comments only",
			dockerfile: `# This is a comment
# Another comment`,
			wantStages: 0,
			wantError:  true,
		},
		{
			name: "complex multi-stage with multiple copy operations",
			dockerfile: `FROM node:18 AS frontend
WORKDIR /app
COPY package.json .
RUN npm install
COPY src/ ./src/
RUN npm run build

FROM golang:1.20 AS backend
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY pkg/ ./pkg/
RUN go build -o server cmd/main.go

FROM alpine:3.18 AS final
RUN apk add --no-cache ca-certificates
COPY --from=frontend /app/dist /static
COPY --from=backend /go/src/app/server /usr/local/bin/server
EXPOSE 8080
CMD ["server"]`,
			wantStages: 3,
			wantError:  false,
			checkStage: func(m *MultiStageBuilder) bool {
				stages := m.GetStages()
				return len(stages) == 3 &&
					stages[0] == "frontend" &&
					stages[1] == "backend" &&
					stages[2] == "final"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewMultiStageBuilder(nil)
			err := builder.ParseDockerfile(tt.dockerfile)

			if (err != nil) != tt.wantError {
				t.Errorf("ParseDockerfile() error = %v, wantError %v", err, tt.wantError)
			}

			if len(builder.stages) != tt.wantStages {
				t.Errorf("got %d stages, want %d", len(builder.stages), tt.wantStages)
			}

			if tt.checkStage != nil && !tt.checkStage(builder) {
				t.Errorf("stage validation failed")
			}
		})
	}
}

func TestMultiStageBuilder_AddStage(t *testing.T) {
	tests := []struct {
		name      string
		stageName string
		baseImage string
		wantError bool
	}{
		{
			name:      "valid stage",
			stageName: "test-stage",
			baseImage: "alpine:3.18",
			wantError: false,
		},
		{
			name:      "empty stage name",
			stageName: "",
			baseImage: "alpine:3.18",
			wantError: true,
		},
		{
			name:      "duplicate stage name",
			stageName: "duplicate",
			baseImage: "alpine:3.18",
			wantError: false, // First addition should succeed
		},
	}

	builder := NewMultiStageBuilder(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := builder.AddStage(tt.stageName, tt.baseImage)

			if (err != nil) != tt.wantError {
				t.Errorf("AddStage() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError && tt.stageName != "" {
				stage, exists := builder.stages[tt.stageName]
				if !exists {
					t.Errorf("stage %s was not added", tt.stageName)
				}
				if stage.BaseImage != tt.baseImage {
					t.Errorf("base image = %v, want %v", stage.BaseImage, tt.baseImage)
				}
			}
		})
	}

	// Test duplicate stage name error
	err := builder.AddStage("duplicate", "alpine:latest")
	if err == nil {
		t.Errorf("expected error for duplicate stage name")
	}
}

func TestMultiStageBuilder_ResolveDependencies(t *testing.T) {
	tests := []struct {
		name         string
		setupBuilder func() *MultiStageBuilder
		wantOrder    []string
		wantError    bool
	}{
		{
			name: "no dependencies",
			setupBuilder: func() *MultiStageBuilder {
				builder := NewMultiStageBuilder(nil)
				builder.AddStage("stage1", "alpine:3.18")
				builder.AddStage("stage2", "alpine:3.18")
				return builder
			},
			wantOrder: []string{"stage1", "stage2"},
			wantError: false,
		},
		{
			name: "linear dependencies",
			setupBuilder: func() *MultiStageBuilder {
				builder := NewMultiStageBuilder(nil)
				dockerfile := `FROM golang:1.20 AS builder
RUN go build
FROM alpine:3.18 AS final
COPY --from=builder /app /app`
				builder.ParseDockerfile(dockerfile)
				return builder
			},
			wantOrder: []string{"builder", "final"},
			wantError: false,
		},
		{
			name: "complex dependencies",
			setupBuilder: func() *MultiStageBuilder {
				builder := NewMultiStageBuilder(nil)
				dockerfile := `FROM node:18 AS frontend
RUN npm build
FROM golang:1.20 AS backend
RUN go build
FROM alpine:3.18 AS final
COPY --from=frontend /dist /static
COPY --from=backend /app /app`
				builder.ParseDockerfile(dockerfile)
				return builder
			},
			wantOrder: []string{"frontend", "backend", "final"},
			wantError: false,
		},
		{
			name: "circular dependency",
			setupBuilder: func() *MultiStageBuilder {
				builder := NewMultiStageBuilder(nil)
				// Manually create circular dependency for testing
				builder.AddStage("stage1", "alpine:3.18")
				builder.AddStage("stage2", "alpine:3.18")
				builder.stages["stage1"].Dependencies = []string{"stage2"}
				builder.stages["stage2"].Dependencies = []string{"stage1"}
				return builder
			},
			wantOrder: nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := tt.setupBuilder()
			order, err := builder.ResolveDependencies()

			if (err != nil) != tt.wantError {
				t.Errorf("ResolveDependencies() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError {
				if len(order) != len(tt.wantOrder) {
					t.Errorf("order length = %v, want %v", len(order), len(tt.wantOrder))
				}

				// Check that all expected stages are present
				for _, expected := range tt.wantOrder {
					if !contains(order, expected) {
						t.Errorf("expected stage %s not found in order", expected)
					}
				}

				// Verify dependency constraints are satisfied
				stagePos := make(map[string]int)
				for i, stage := range order {
					stagePos[stage] = i
				}

				for _, stage := range order {
					stageObj := builder.stages[stage]
					for _, dep := range stageObj.Dependencies {
						if stagePos[dep] >= stagePos[stage] {
							t.Errorf("dependency %s comes after %s in order", dep, stage)
						}
					}
				}
			}
		})
	}
}

func TestMultiStageBuilder_HandleCopyFromStage(t *testing.T) {
	tests := []struct {
		name        string
		setupStages func(*MultiStageBuilder)
		sourceStage string
		sourcePath  string
		destPath    string
		wantError   bool
	}{
		{
			name: "valid copy operation",
			setupStages: func(b *MultiStageBuilder) {
				b.AddStage("builder", "golang:1.20")
				b.AddStage("final", "alpine:3.18")
				b.currentStage = "final"
			},
			sourceStage: "builder",
			sourcePath:  "/app/binary",
			destPath:    "/usr/local/bin/app",
			wantError:   false,
		},
		{
			name: "nonexistent source stage",
			setupStages: func(b *MultiStageBuilder) {
				b.AddStage("final", "alpine:3.18")
				b.currentStage = "final"
			},
			sourceStage: "nonexistent",
			sourcePath:  "/app/binary",
			destPath:    "/usr/local/bin/app",
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtx := NewMockBuildContextManager()
			builder := NewMultiStageBuilder(mockCtx)
			tt.setupStages(builder)

			err := builder.HandleCopyFromStage(tt.sourceStage, tt.sourcePath, tt.destPath)

			if (err != nil) != tt.wantError {
				t.Errorf("HandleCopyFromStage() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError && builder.currentStage != "" {
				stage := builder.stages[builder.currentStage]
				expectedArtifact := fmt.Sprintf("%s:%s", tt.sourceStage, tt.sourcePath)
				if stage.Artifacts[tt.destPath] != expectedArtifact {
					t.Errorf("artifact not recorded correctly, got %v, want %v",
						stage.Artifacts[tt.destPath], expectedArtifact)
				}
			}
		})
	}
}

func TestMultiStageBuilder_ProcessStage(t *testing.T) {
	tests := []struct {
		name      string
		stageName string
		wantError bool
	}{
		{
			name:      "valid stage",
			stageName: "test-stage",
			wantError: false,
		},
		{
			name:      "nonexistent stage",
			stageName: "nonexistent",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtx := NewMockBuildContextManager()
			builder := NewMultiStageBuilder(mockCtx)

			if tt.stageName == "test-stage" {
				builder.AddStage("test-stage", "alpine:3.18")
			}

			err := builder.ProcessStage(tt.stageName)

			if (err != nil) != tt.wantError {
				t.Errorf("ProcessStage() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError {
				if builder.currentStage != tt.stageName {
					t.Errorf("current stage = %v, want %v", builder.currentStage, tt.stageName)
				}
			}
		})
	}
}

func TestMultiStageBuilder_GetStageArtifacts(t *testing.T) {
	mockCtx := NewMockBuildContextManager()
	builder := NewMultiStageBuilder(mockCtx)
	builder.AddStage("test-stage", "alpine:3.18")

	// Test getting artifacts from existing stage
	artifacts, err := builder.GetStageArtifacts("test-stage")
	if err != nil {
		t.Errorf("GetStageArtifacts() error = %v", err)
	}
	if artifacts == nil {
		t.Errorf("GetStageArtifacts() returned nil artifacts")
	}

	// Test getting artifacts from nonexistent stage
	_, err = builder.GetStageArtifacts("nonexistent")
	if err == nil {
		t.Errorf("expected error for nonexistent stage")
	}
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}