package registry

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/oci"
)

type MockRegistry struct {
	mu           sync.RWMutex
	repositories map[string]*MockRepository
	errors       map[string]error
	delays       map[string]time.Duration
	callCounts   map[string]int
}

type MockRepository struct {
	Name   string
	Images map[string]*MockImage
}

type MockImage struct {
	Tag    string
	Digest string
	Size   int64
	Data   []byte
}

func NewMockRegistry() *MockRegistry {
	return &MockRegistry{
		repositories: make(map[string]*MockRepository),
		errors:       make(map[string]error),
		delays:       make(map[string]time.Duration),
		callCounts:   make(map[string]int),
	}
}

func (m *MockRegistry) InjectError(operation string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errors[operation] = err
}

func (m *MockRegistry) InjectDelay(operation string, delay time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.delays[operation] = delay
}

func (m *MockRegistry) GetCallCount(operation string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.callCounts[operation]
}

func (m *MockRegistry) checkError(operation string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.callCounts[operation]++

	if delay, exists := m.delays[operation]; exists {
		time.Sleep(delay)
	}

	if err, exists := m.errors[operation]; exists {
		return err
	}

	return nil
}

func (m *MockRegistry) Push(ctx context.Context, image string, content io.Reader) error {
	if err := m.checkError("Push"); err != nil {
		return err
	}

	repository, tag, err := ParseImageRef(image)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.repositories[repository]; !exists {
		m.repositories[repository] = &MockRepository{Name: repository, Images: make(map[string]*MockImage)}
	}

	data, err := io.ReadAll(content)
	if err != nil {
		return err
	}

	m.repositories[repository].Images[tag] = &MockImage{Tag: tag, Digest: calculateDigest(data), Size: int64(len(data)), Data: data}
	return nil
}

func (m *MockRegistry) List(ctx context.Context) ([]string, error) {
	if err := m.checkError("List"); err != nil {
		return nil, err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	repositories := make([]string, 0, len(m.repositories))
	for name := range m.repositories {
		repositories = append(repositories, name)
	}
	return repositories, nil
}

func (m *MockRegistry) Exists(ctx context.Context, repository string) (bool, error) {
	if err := m.checkError("Exists"); err != nil {
		return false, err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.repositories[repository]
	return exists, nil
}

func (m *MockRegistry) Delete(ctx context.Context, repository string) error {
	if err := m.checkError("Delete"); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.repositories[repository]; !exists {
		return fmt.Errorf("repository %s not found", repository)
	}
	delete(m.repositories, repository)
	return nil
}

func (m *MockRegistry) Close() error {
	return m.checkError("Close")
}

func (m *MockRegistry) AddRepository(name string) *MockRepository {
	m.mu.Lock()
	defer m.mu.Unlock()
	repo := &MockRepository{Name: name, Images: make(map[string]*MockImage)}
	m.repositories[name] = repo
	return repo
}

func (m *MockRegistry) AddImage(repository, tag string, size int64) *MockImage {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.repositories[repository]; !exists {
		m.repositories[repository] = &MockRepository{Name: repository, Images: make(map[string]*MockImage)}
	}

	data := bytes.Repeat([]byte("mock"), int(size/4)+1)[:size]
	image := &MockImage{Tag: tag, Digest: calculateDigest(data), Size: size, Data: data}
	m.repositories[repository].Images[tag] = image
	return image
}

func (m *MockRegistry) ClearAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.repositories = make(map[string]*MockRepository)
	m.errors = make(map[string]error)
	m.delays = make(map[string]time.Duration)
	m.callCounts = make(map[string]int)
}

type TestHelper struct{ mock *MockRegistry }

func NewTestHelper() *TestHelper { return &TestHelper{NewMockRegistry()} }

func (th *TestHelper) GetMockRegistry() *MockRegistry { return th.mock }

func (th *TestHelper) SetupBasicScenario() {
	th.mock.AddRepository("nginx")
	th.mock.AddImage("nginx", "latest", 50*1024*1024)
}

func (th *TestHelper) CreateTestManifest() *oci.Manifest {
	return &oci.Manifest{
		SchemaVersion: 2,
		MediaType:     "application/vnd.docker.distribution.manifest.v2+json",
		Config:        oci.Descriptor{Digest: "sha256:cfg", Size: 1024, MediaType: "application/vnd.docker.container.image.v1+json"},
		Layers:        []oci.Descriptor{{Digest: "sha256:layer1", Size: 2048, MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip"}},
	}
}

// ParseImageRef parses an image reference into repository and tag
func ParseImageRef(image string) (string, string, error) {
	parts := strings.Split(image, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid image reference: %s", image)
	}
	return parts[0], parts[1], nil
}

// calculateDigest calculates SHA256 digest of data
func calculateDigest(data []byte) string {
	h := sha256.Sum256(data)
	return fmt.Sprintf("sha256:%x", h)
}
