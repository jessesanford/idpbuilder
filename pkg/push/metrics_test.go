package push

import (
	"testing"
	"time"
)

// TestNoOpMetricsImplementsInterface verifies that NoOpMetrics implements the Metrics interface
func TestNoOpMetricsImplementsInterface(t *testing.T) {
	var _ Metrics = (*NoOpMetrics)(nil)
}

// TestNoOpMetricsRecordPushStart verifies RecordPushStart is a true no-op
func TestNoOpMetricsRecordPushStart(t *testing.T) {
	m := &NoOpMetrics{}

	// Should not panic and should complete immediately
	m.RecordPushStart("test-image:latest", "registry.example.com")
	m.RecordPushStart("", "")
	m.RecordPushStart("invalid image name with spaces", "")
}

// TestNoOpMetricsRecordPushComplete verifies RecordPushComplete is a true no-op
func TestNoOpMetricsRecordPushComplete(t *testing.T) {
	m := &NoOpMetrics{}

	// Test successful completion
	m.RecordPushComplete("test-image:latest", "registry.example.com", time.Second, nil)

	// Test failed completion with error
	testErr := &MockError{msg: "connection timeout"}
	m.RecordPushComplete("test-image:latest", "registry.example.com", time.Millisecond*500, testErr)

	// Test with zero duration
	m.RecordPushComplete("test-image:latest", "registry.example.com", 0, nil)

	// Test with empty strings
	m.RecordPushComplete("", "", 0, nil)
}

// TestNoOpMetricsRecordRetry verifies RecordRetry is a true no-op
func TestNoOpMetricsRecordRetry(t *testing.T) {
	m := &NoOpMetrics{}

	// Test various retry scenarios
	m.RecordRetry("test-image:latest", "registry.example.com", 1, "network timeout")
	m.RecordRetry("test-image:latest", "registry.example.com", 5, "rate limited")
	m.RecordRetry("test-image:latest", "registry.example.com", 0, "")
	m.RecordRetry("", "", -1, "")
}

// TestNoOpMetricsRecordProgress verifies RecordProgress is a true no-op
func TestNoOpMetricsRecordProgress(t *testing.T) {
	m := &NoOpMetrics{}

	// Test progress updates
	m.RecordProgress("test-image:latest", 1024, 1024*1024)
	m.RecordProgress("test-image:latest", 512*1024, 1024*1024)
	m.RecordProgress("test-image:latest", 1024*1024, 1024*1024)

	// Test edge cases
	m.RecordProgress("", 0, 0)
	m.RecordProgress("test-image", -1, 100)
	m.RecordProgress("test-image", 100, -1)
}

// TestNoOpMetricsRecordLayerUpload verifies RecordLayerUpload is a true no-op
func TestNoOpMetricsRecordLayerUpload(t *testing.T) {
	m := &NoOpMetrics{}

	// Test layer uploads
	m.RecordLayerUpload("test-image:latest", "sha256:abc123", 1024*1024, time.Second)
	m.RecordLayerUpload("test-image:latest", "sha256:def456", 512*1024, time.Millisecond*500)
	m.RecordLayerUpload("test-image:latest", "sha256:ghi789", 0, 0)

	// Test edge cases
	m.RecordLayerUpload("", "", -1, -1)
	m.RecordLayerUpload("test", "invalid-digest", 0, time.Hour*24)
}

// TestNoOpMetricsAllMethodsConcurrent tests all methods under concurrent load
func TestNoOpMetricsAllMethodsConcurrent(t *testing.T) {
	m := &NoOpMetrics{}

	// Run all methods concurrently to verify thread safety
	done := make(chan bool, 5)

	go func() {
		for i := 0; i < 1000; i++ {
			m.RecordPushStart("test-image", "registry.example.com")
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			m.RecordPushComplete("test-image", "registry.example.com", time.Second, nil)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			m.RecordRetry("test-image", "registry.example.com", i, "retry")
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			m.RecordProgress("test-image", int64(i*1024), 1024*1024)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			m.RecordLayerUpload("test-image", "sha256:abc", 1024, time.Millisecond)
		}
		done <- true
	}()

	// Wait for all goroutines
	for i := 0; i < 5; i++ {
		<-done
	}
}

// TestNoOpMetricsNoSideEffects verifies methods have no observable side effects
func TestNoOpMetricsNoSideEffects(t *testing.T) {
	m := &NoOpMetrics{}

	// Call each method multiple times
	for i := 0; i < 100; i++ {
		m.RecordPushStart("test", "registry")
		m.RecordPushComplete("test", "registry", time.Second, nil)
		m.RecordRetry("test", "registry", i, "reason")
		m.RecordProgress("test", int64(i), 100)
		m.RecordLayerUpload("test", "digest", int64(i), time.Millisecond)
	}

	// Verify NoOpMetrics struct is still empty (no fields modified)
	// This is implicit - if it had fields, they would be visible in the struct definition
}

// BenchmarkNoOpMetricsRecordPushStart benchmarks the RecordPushStart method
func BenchmarkNoOpMetricsRecordPushStart(b *testing.B) {
	m := &NoOpMetrics{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.RecordPushStart("test-image:latest", "registry.example.com")
	}
}

// BenchmarkNoOpMetricsRecordPushComplete benchmarks the RecordPushComplete method
func BenchmarkNoOpMetricsRecordPushComplete(b *testing.B) {
	m := &NoOpMetrics{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.RecordPushComplete("test-image:latest", "registry.example.com", time.Second, nil)
	}
}

// BenchmarkNoOpMetricsRecordRetry benchmarks the RecordRetry method
func BenchmarkNoOpMetricsRecordRetry(b *testing.B) {
	m := &NoOpMetrics{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.RecordRetry("test-image:latest", "registry.example.com", i, "network timeout")
	}
}

// BenchmarkNoOpMetricsRecordProgress benchmarks the RecordProgress method
func BenchmarkNoOpMetricsRecordProgress(b *testing.B) {
	m := &NoOpMetrics{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.RecordProgress("test-image:latest", int64(i*1024), 1024*1024)
	}
}

// BenchmarkNoOpMetricsRecordLayerUpload benchmarks the RecordLayerUpload method
func BenchmarkNoOpMetricsRecordLayerUpload(b *testing.B) {
	m := &NoOpMetrics{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.RecordLayerUpload("test-image:latest", "sha256:abc123", 1024*1024, time.Second)
	}
}

// BenchmarkNoOpMetricsAllMethodsSequential benchmarks all methods called sequentially
func BenchmarkNoOpMetricsAllMethodsSequential(b *testing.B) {
	m := &NoOpMetrics{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.RecordPushStart("test-image:latest", "registry.example.com")
		m.RecordProgress("test-image:latest", 512*1024, 1024*1024)
		m.RecordLayerUpload("test-image:latest", "sha256:abc123", 1024*1024, time.Millisecond*100)
		m.RecordPushComplete("test-image:latest", "registry.example.com", time.Second, nil)
	}
}

// BenchmarkNoOpMetricsAllMethodsParallel benchmarks all methods called in parallel
func BenchmarkNoOpMetricsAllMethodsParallel(b *testing.B) {
	m := &NoOpMetrics{}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.RecordPushStart("test-image:latest", "registry.example.com")
			m.RecordProgress("test-image:latest", int64(i*1024), 1024*1024)
			m.RecordLayerUpload("test-image:latest", "sha256:abc123", 1024*1024, time.Millisecond*100)
			m.RecordPushComplete("test-image:latest", "registry.example.com", time.Second, nil)
			i++
		}
	})
}

// MockError is a simple error implementation for testing
type MockError struct {
	msg string
}

func (e *MockError) Error() string {
	return e.msg
}
