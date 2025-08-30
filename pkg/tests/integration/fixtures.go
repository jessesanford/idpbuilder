package integration

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"time"
)

// GetTestDockerfile returns the path to a test Dockerfile
func GetTestDockerfile(name string) string {
	return filepath.Join("test-data", "dockerfiles", fmt.Sprintf("%s.Dockerfile", name))
}

// GetTestContext returns the path to a test build context
func GetTestContext(name string) string {
	return filepath.Join("test-data", "contexts", name)
}

// GenerateTestImageTag generates a unique image tag for testing
func GenerateTestImageTag() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("integration-test-%d:latest", rand.Int31())
}

// GenerateGiteaImageTag generates a test image tag for Gitea registry
func GenerateGiteaImageTag() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("gitea.local:443/test/integration-%d:v1", rand.Int31())
}

// TestDockerfiles contains the content for test Dockerfiles
var TestDockerfiles = map[string]string{
	"simple": `FROM alpine:latest
LABEL maintainer="integration-test"
CMD ["echo", "Hello from integration test"]
`,
	"multistage": `FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY main.go .
RUN go build -o app main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/app .
CMD ["./app"]
`,
	"invalid": `FROM alpine:latest
RUN invalid-command-that-does-not-exist
INVALID SYNTAX HERE
CMD ["echo", "this will fail"]
`,
}

// TestContextFiles contains test application files
var TestContextFiles = map[string]string{
	"simple-app/main.go": `package main

import "fmt"

func main() {
	fmt.Println("Integration test application")
	fmt.Println("CLI build and push test successful!")
}
`,
}