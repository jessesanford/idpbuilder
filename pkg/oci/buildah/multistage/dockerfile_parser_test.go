package multistage

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDockerfileParser_Parse(t *testing.T) {
	tests := []struct {
		name           string
		dockerfile     string
		expectedStages int
		expectedOrder  []string
		expectError    bool
	}{
		{
			name: "simple_multi_stage",
			dockerfile: `FROM node:16 AS builder
RUN npm install
COPY . .
RUN npm run build

FROM nginx:alpine AS production  
COPY --from=builder /app/dist /usr/share/nginx/html`,
			expectedStages: 2,
			expectedOrder:  []string{"builder", "production"},
			expectError:    false,
		},
		{
			name: "complex_dependency_graph",
			dockerfile: `FROM golang:1.19 AS base
RUN go version

FROM base AS deps
COPY go.mod go.sum ./
RUN go mod download

FROM base AS build
COPY --from=deps /go/pkg /go/pkg
COPY . .
RUN go build -o app

FROM alpine:latest AS final
COPY --from=build /app /app`,
			expectedStages: 4,
			expectedOrder:  []string{"base", "deps", "build", "final"},
			expectError:    false,
		},
		{
			name: "unnamed_stages",
			dockerfile: `FROM node:16
RUN npm install

FROM nginx:alpine
COPY --from=stage-0 /app /usr/share/nginx/html`,
			expectedStages: 2,
			expectedOrder:  []string{"stage-0", "stage-1"},
			expectError:    false,
		},
		{
			name: "circular_dependency",
			dockerfile: `FROM alpine AS stage1
COPY --from=stage2 /file1 .

FROM stage1 AS stage2
COPY --from=stage1 /file2 .`,
			expectedStages: 2,
			expectError:    true,
		},
		{
			name: "undefined_stage_reference",
			dockerfile: `FROM alpine AS stage1
COPY --from=nonexistent /file .`,
			expectedStages: 1,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewDockerfileParser()
			graph, err := parser.Parse(strings.NewReader(tt.dockerfile))

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, graph.Stages, tt.expectedStages)
			assert.Equal(t, tt.expectedOrder, graph.ExecutionOrder)
		})
	}
}

func TestDockerfileParser_parseCommand(t *testing.T) {
	parser := NewDockerfileParser()
	tests := []struct {
		name        string
		line        string
		expectedCmd Command
		expectError bool
	}{
		{
			name: "simple_run",
			line: "RUN npm install",
			expectedCmd: Command{
				Type: "RUN",
				Args: []string{"npm", "install"},
			},
			expectError: false,
		},
		{
			name: "copy_with_from",
			line: "COPY --from=builder /app/dist /usr/share/nginx/html",
			expectedCmd: Command{
				Type: "COPY",
				From: "builder",
				Args: []string{"/app/dist", "/usr/share/nginx/html"},
			},
			expectError: false,
		},
		{
			name: "add_command",
			line: "ADD https://example.com/file.tar.gz /tmp/",
			expectedCmd: Command{
				Type: "ADD",
				Args: []string{"https://example.com/file.tar.gz", "/tmp/"},
			},
			expectError: false,
		},
		{
			name:        "empty_command",
			line:        "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := parser.parseCommand(tt.line)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedCmd, cmd)
		})
	}
}

func TestDockerfileParser_calculateExecutionOrder(t *testing.T) {
	parser := NewDockerfileParser()

	tests := []struct {
		name         string
		stages       []BuildStage
		dependencies map[string][]string
		expectedOrder []string
		expectError   bool
	}{
		{
			name: "linear_dependency",
			stages: []BuildStage{
				{Name: "stage1"},
				{Name: "stage2"},
				{Name: "stage3"},
			},
			dependencies: map[string][]string{
				"stage1": {},
				"stage2": {"stage1"},
				"stage3": {"stage2"},
			},
			expectedOrder: []string{"stage1", "stage2", "stage3"},
			expectError:   false,
		},
		{
			name: "parallel_stages",
			stages: []BuildStage{
				{Name: "base"},
				{Name: "frontend"},
				{Name: "backend"},
				{Name: "final"},
			},
			dependencies: map[string][]string{
				"base":     {},
				"frontend": {"base"},
				"backend":  {"base"},
				"final":    {"frontend", "backend"},
			},
			expectedOrder: []string{"base", "frontend", "backend", "final"},
			expectError:   false,
		},
		{
			name: "circular_dependency",
			stages: []BuildStage{
				{Name: "stage1"},
				{Name: "stage2"},
			},
			dependencies: map[string][]string{
				"stage1": {"stage2"},
				"stage2": {"stage1"},
			},
			expectError: true,
		},
		{
			name: "undefined_dependency",
			stages: []BuildStage{
				{Name: "stage1"},
			},
			dependencies: map[string][]string{
				"stage1": {"undefined"},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := parser.calculateExecutionOrder(tt.stages, tt.dependencies)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedOrder, order)
		})
	}
}

func TestDockerfileParser_stagePatternMatching(t *testing.T) {
	parser := NewDockerfileParser()
	tests := []struct {
		line          string
		shouldMatch   bool
		expectedImage string
		expectedStage string
	}{
		{"FROM node:16 AS builder", true, "node:16", "builder"},
		{"FROM nginx:alpine", true, "nginx:alpine", ""},
		{"from ubuntu:20.04 as base", true, "ubuntu:20.04", "base"},
		{"RUN echo hello", false, "", ""},
		{"# Comment with FROM", false, "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			matches := parser.stagePattern.FindStringSubmatch(tt.line)
			if !tt.shouldMatch {
				assert.Nil(t, matches)
				return
			}

			require.NotNil(t, matches)
			assert.Equal(t, tt.expectedImage, matches[1])
			assert.Equal(t, tt.expectedStage, matches[2])
		})
	}
}