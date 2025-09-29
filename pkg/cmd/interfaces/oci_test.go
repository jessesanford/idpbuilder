package interfaces

import (
	"context"
	"testing"
)

// Mock implementations for testing interface compliance

// mockOCICommand implements OCICommand for testing
type mockOCICommand struct {
	executed    bool
	validated   bool
	description string
}

func (m *mockOCICommand) Execute(ctx context.Context, opts CommandOptions) error {
	m.executed = true
	return nil
}

func (m *mockOCICommand) Validate(opts CommandOptions) error {
	m.validated = true
	return nil
}

func (m *mockOCICommand) GetDescription() string {
	return m.description
}

// mockPushCommand implements PushCommand for testing
type mockPushCommand struct {
	mockOCICommand
	progressReporter ProgressReporter
	pushExecuted     bool
	pushValidated    bool
	capabilities     []string
}

func (m *mockPushCommand) SetProgressReporter(reporter ProgressReporter) {
	m.progressReporter = reporter
}

func (m *mockPushCommand) ExecutePush(ctx context.Context, opts PushOptions) error {
	m.pushExecuted = true
	return nil
}

func (m *mockPushCommand) ValidatePush(opts PushOptions) error {
	m.pushValidated = true
	return nil
}

func (m *mockPushCommand) GetPushCapabilities() []string {
	return m.capabilities
}

// mockPullCommand implements PullCommand for testing
type mockPullCommand struct {
	mockOCICommand
	progressReporter ProgressReporter
	extractor        Extractor
	pullExecuted     bool
	pullValidated    bool
	capabilities     []string
}

func (m *mockPullCommand) SetProgressReporter(reporter ProgressReporter) {
	m.progressReporter = reporter
}

func (m *mockPullCommand) SetExtractor(extractor Extractor) {
	m.extractor = extractor
}

func (m *mockPullCommand) ExecutePull(ctx context.Context, opts PullOptions) error {
	m.pullExecuted = true
	return nil
}

func (m *mockPullCommand) ValidatePull(opts PullOptions) error {
	m.pullValidated = true
	return nil
}

func (m *mockPullCommand) GetPullCapabilities() []string {
	return m.capabilities
}

// mockListCommand implements ListCommand for testing
type mockListCommand struct {
	mockOCICommand
	listExecuted    bool
	listValidated   bool
	capabilities    []string
	supportsPaging  bool
}

func (m *mockListCommand) ExecuteList(ctx context.Context, opts ListOptions) error {
	m.listExecuted = true
	return nil
}

func (m *mockListCommand) ValidateList(opts ListOptions) error {
	m.listValidated = true
	return nil
}

func (m *mockListCommand) GetListCapabilities() []string {
	return m.capabilities
}

func (m *mockListCommand) SupportsPagination() bool {
	return m.supportsPaging
}

// mockInspectCommand implements InspectCommand for testing
type mockInspectCommand struct {
	mockOCICommand
	inspectExecuted bool
	inspectValidated bool
	capabilities    []string
	formats         []string
}

func (m *mockInspectCommand) ExecuteInspect(ctx context.Context, opts InspectOptions) error {
	m.inspectExecuted = true
	return nil
}

func (m *mockInspectCommand) ValidateInspect(opts InspectOptions) error {
	m.inspectValidated = true
	return nil
}

func (m *mockInspectCommand) GetInspectCapabilities() []string {
	return m.capabilities
}

func (m *mockInspectCommand) GetSupportedFormats() []string {
	return m.formats
}

// mockProgressReporter implements ProgressReporter for testing
type mockProgressReporter struct {
	started    bool
	completed  bool
	errored    bool
	progress   []string
}

func (m *mockProgressReporter) ReportProgress(current, total int64, message string) {
	m.progress = append(m.progress, message)
}

func (m *mockProgressReporter) Start(message string) {
	m.started = true
}

func (m *mockProgressReporter) Complete(message string) {
	m.completed = true
}

func (m *mockProgressReporter) Error(err error) {
	m.errored = true
}

// mockExtractor implements Extractor for testing
type mockExtractor struct {
	extracted bool
	formats   []string
}

func (m *mockExtractor) Extract(source, destination string) error {
	m.extracted = true
	return nil
}

func (m *mockExtractor) SupportedFormats() []string {
	return m.formats
}

// Test OCICommand interface compliance
func TestOCICommandInterface(t *testing.T) {
	cmd := &mockOCICommand{description: "test command"}

	// Test interface compliance
	var _ OCICommand = cmd

	ctx := context.Background()
	opts := CommandOptions{Registry: "test.registry.com"}

	// Test Execute
	err := cmd.Execute(ctx, opts)
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}
	if !cmd.executed {
		t.Error("Execute was not called")
	}

	// Test Validate
	err = cmd.Validate(opts)
	if err != nil {
		t.Errorf("Validate failed: %v", err)
	}
	if !cmd.validated {
		t.Error("Validate was not called")
	}

	// Test GetDescription
	desc := cmd.GetDescription()
	if desc != "test command" {
		t.Errorf("Expected 'test command', got '%s'", desc)
	}
}

// Test PushCommand interface compliance and composition
func TestPushCommandInterface(t *testing.T) {
	cmd := &mockPushCommand{
		mockOCICommand: mockOCICommand{description: "push command"},
		capabilities:   []string{"push", "force-push"},
	}

	// Test interface compliance
	var _ PushCommand = cmd
	var _ OCICommand = cmd // Should also implement OCICommand

	ctx := context.Background()
	pushOpts := PushOptions{
		CommandOptions: CommandOptions{Registry: "test.registry.com"},
		Source:         "/local/path",
		Destination:    "test.registry.com/repo:tag",
	}

	// Test PushCommand methods
	err := cmd.ExecutePush(ctx, pushOpts)
	if err != nil {
		t.Errorf("ExecutePush failed: %v", err)
	}
	if !cmd.pushExecuted {
		t.Error("ExecutePush was not called")
	}

	err = cmd.ValidatePush(pushOpts)
	if err != nil {
		t.Errorf("ValidatePush failed: %v", err)
	}
	if !cmd.pushValidated {
		t.Error("ValidatePush was not called")
	}

	caps := cmd.GetPushCapabilities()
	if len(caps) != 2 {
		t.Errorf("Expected 2 capabilities, got %d", len(caps))
	}

	// Test progress reporter
	reporter := &mockProgressReporter{}
	cmd.SetProgressReporter(reporter)
	if cmd.progressReporter != reporter {
		t.Error("SetProgressReporter did not set the reporter")
	}

	// Test inherited OCICommand methods
	err = cmd.Execute(ctx, pushOpts.CommandOptions)
	if err != nil {
		t.Errorf("Inherited Execute failed: %v", err)
	}
	if !cmd.executed {
		t.Error("Inherited Execute was not called")
	}
}

// Test PullCommand interface compliance and composition
func TestPullCommandInterface(t *testing.T) {
	cmd := &mockPullCommand{
		mockOCICommand: mockOCICommand{description: "pull command"},
		capabilities:   []string{"pull", "extract"},
	}

	// Test interface compliance
	var _ PullCommand = cmd
	var _ OCICommand = cmd // Should also implement OCICommand

	ctx := context.Background()
	pullOpts := PullOptions{
		CommandOptions: CommandOptions{Registry: "test.registry.com"},
		Source:         "test.registry.com/repo:tag",
		Destination:    "/local/path",
	}

	// Test PullCommand methods
	err := cmd.ExecutePull(ctx, pullOpts)
	if err != nil {
		t.Errorf("ExecutePull failed: %v", err)
	}
	if !cmd.pullExecuted {
		t.Error("ExecutePull was not called")
	}

	err = cmd.ValidatePull(pullOpts)
	if err != nil {
		t.Errorf("ValidatePull failed: %v", err)
	}
	if !cmd.pullValidated {
		t.Error("ValidatePull was not called")
	}

	caps := cmd.GetPullCapabilities()
	if len(caps) != 2 {
		t.Errorf("Expected 2 capabilities, got %d", len(caps))
	}

	// Test progress reporter
	reporter := &mockProgressReporter{}
	cmd.SetProgressReporter(reporter)
	if cmd.progressReporter != reporter {
		t.Error("SetProgressReporter did not set the reporter")
	}

	// Test extractor
	extractor := &mockExtractor{formats: []string{"tar", "zip"}}
	cmd.SetExtractor(extractor)
	if cmd.extractor != extractor {
		t.Error("SetExtractor did not set the extractor")
	}
}

// Test ListCommand interface compliance
func TestListCommandInterface(t *testing.T) {
	cmd := &mockListCommand{
		mockOCICommand: mockOCICommand{description: "list command"},
		capabilities:   []string{"list", "search"},
		supportsPaging: true,
	}

	// Test interface compliance
	var _ ListCommand = cmd
	var _ OCICommand = cmd

	ctx := context.Background()
	listOpts := ListOptions{
		CommandOptions: CommandOptions{Registry: "test.registry.com"},
		Repository:     "test/repo",
		Tags:           true,
		Limit:          10,
	}

	// Test ListCommand methods
	err := cmd.ExecuteList(ctx, listOpts)
	if err != nil {
		t.Errorf("ExecuteList failed: %v", err)
	}
	if !cmd.listExecuted {
		t.Error("ExecuteList was not called")
	}

	err = cmd.ValidateList(listOpts)
	if err != nil {
		t.Errorf("ValidateList failed: %v", err)
	}
	if !cmd.listValidated {
		t.Error("ValidateList was not called")
	}

	caps := cmd.GetListCapabilities()
	if len(caps) != 2 {
		t.Errorf("Expected 2 capabilities, got %d", len(caps))
	}

	if !cmd.SupportsPagination() {
		t.Error("Expected pagination support")
	}
}

// Test InspectCommand interface compliance
func TestInspectCommandInterface(t *testing.T) {
	cmd := &mockInspectCommand{
		mockOCICommand: mockOCICommand{description: "inspect command"},
		capabilities:   []string{"inspect", "raw"},
		formats:        []string{"json", "yaml", "table"},
	}

	// Test interface compliance
	var _ InspectCommand = cmd
	var _ OCICommand = cmd

	ctx := context.Background()
	inspectOpts := InspectOptions{
		CommandOptions: CommandOptions{Registry: "test.registry.com"},
		Reference:      "test.registry.com/repo:tag",
		Raw:            false,
		Format:         "json",
	}

	// Test InspectCommand methods
	err := cmd.ExecuteInspect(ctx, inspectOpts)
	if err != nil {
		t.Errorf("ExecuteInspect failed: %v", err)
	}
	if !cmd.inspectExecuted {
		t.Error("ExecuteInspect was not called")
	}

	err = cmd.ValidateInspect(inspectOpts)
	if err != nil {
		t.Errorf("ValidateInspect failed: %v", err)
	}
	if !cmd.inspectValidated {
		t.Error("ValidateInspect was not called")
	}

	caps := cmd.GetInspectCapabilities()
	if len(caps) != 2 {
		t.Errorf("Expected 2 capabilities, got %d", len(caps))
	}

	formats := cmd.GetSupportedFormats()
	if len(formats) != 3 {
		t.Errorf("Expected 3 formats, got %d", len(formats))
	}
}

// Test ProgressReporter interface compliance
func TestProgressReporterInterface(t *testing.T) {
	reporter := &mockProgressReporter{}

	// Test interface compliance
	var _ ProgressReporter = reporter

	// Test methods
	reporter.Start("starting")
	if !reporter.started {
		t.Error("Start was not called")
	}

	reporter.ReportProgress(50, 100, "progress message")
	if len(reporter.progress) != 1 {
		t.Error("ReportProgress was not called")
	}

	reporter.Complete("completed")
	if !reporter.completed {
		t.Error("Complete was not called")
	}

	reporter.Error(nil)
	if !reporter.errored {
		t.Error("Error was not called")
	}
}

// Test Extractor interface compliance
func TestExtractorInterface(t *testing.T) {
	extractor := &mockExtractor{formats: []string{"tar", "zip"}}

	// Test interface compliance
	var _ Extractor = extractor

	// Test methods
	err := extractor.Extract("source", "dest")
	if err != nil {
		t.Errorf("Extract failed: %v", err)
	}
	if !extractor.extracted {
		t.Error("Extract was not called")
	}

	formats := extractor.SupportedFormats()
	if len(formats) != 2 {
		t.Errorf("Expected 2 formats, got %d", len(formats))
	}
}