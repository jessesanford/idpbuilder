package push

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"
)

// TestStreamingPusherBufferPool verifies buffer pool operations
func TestStreamingPusherBufferPool(t *testing.T) {
	sp := NewStreamingPusher()

	// Test buffer acquisition
	buf1 := sp.GetBuffer()
	if buf1 == nil {
		t.Fatal("GetBuffer returned nil")
	}
	if len(*buf1) != DefaultBufferSize {
		t.Errorf("Buffer size mismatch: got %d, want %d", len(*buf1), DefaultBufferSize)
	}

	// Test buffer return and reacquisition
	sp.PutBuffer(buf1)
	buf2 := sp.GetBuffer()
	if buf2 == nil {
		t.Fatal("GetBuffer returned nil after PutBuffer")
	}

	// Verify buffer is reused (same pointer)
	if buf1 != buf2 {
		t.Log("Buffer was not reused (pool created new buffer, which is acceptable)")
	}

	sp.PutBuffer(buf2)
}

// TestStreamingPusherOptions verifies functional options
func TestStreamingPusherOptions(t *testing.T) {
	tests := []struct {
		name     string
		opts     []StreamingOption
		wantSize int
		wantMax  int
	}{
		{
			name:     "default options",
			opts:     nil,
			wantSize: DefaultChunkSize,
			wantMax:  DefaultMaxConcurrentLayers,
		},
		{
			name:     "custom chunk size",
			opts:     []StreamingOption{WithChunkSize(2 * 1024 * 1024)},
			wantSize: 2 * 1024 * 1024,
			wantMax:  DefaultMaxConcurrentLayers,
		},
		{
			name:     "custom max concurrent",
			opts:     []StreamingOption{WithMaxConcurrentOps(8)},
			wantSize: DefaultChunkSize,
			wantMax:  8,
		},
		{
			name: "both options",
			opts: []StreamingOption{
				WithChunkSize(512 * 1024),
				WithMaxConcurrentOps(16),
			},
			wantSize: 512 * 1024,
			wantMax:  16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := NewStreamingPusher(tt.opts...)

			if sp.chunkSize != tt.wantSize {
				t.Errorf("chunkSize = %d, want %d", sp.chunkSize, tt.wantSize)
			}
			if sp.maxConcurrentOps != tt.wantMax {
				t.Errorf("maxConcurrentOps = %d, want %d", sp.maxConcurrentOps, tt.wantMax)
			}
			if cap(sp.semaphore) != tt.wantMax {
				t.Errorf("semaphore capacity = %d, want %d", cap(sp.semaphore), tt.wantMax)
			}
		})
	}
}

// TestStreamingPusherConcurrency verifies semaphore-based concurrency control
func TestStreamingPusherConcurrency(t *testing.T) {
	sp := NewStreamingPusher(WithMaxConcurrentOps(2))
	ctx := context.Background()

	// Acquire first slot
	if err := sp.AcquireSlot(ctx); err != nil {
		t.Fatalf("Failed to acquire first slot: %v", err)
	}

	// Acquire second slot
	if err := sp.AcquireSlot(ctx); err != nil {
		t.Fatalf("Failed to acquire second slot: %v", err)
	}

	// Third acquire should block - test with timeout
	ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := sp.AcquireSlot(ctx2)
	if err != context.DeadlineExceeded {
		t.Errorf("Expected context.DeadlineExceeded, got: %v", err)
	}

	// Release one slot
	sp.ReleaseSlot()

	// Now third acquire should succeed
	if err := sp.AcquireSlot(ctx); err != nil {
		t.Fatalf("Failed to acquire third slot after release: %v", err)
	}

	// Cleanup
	sp.ReleaseSlot()
	sp.ReleaseSlot()
}

// TestStreamingPusherConcurrencyWithCancellation verifies context cancellation
func TestStreamingPusherConcurrencyWithCancellation(t *testing.T) {
	sp := NewStreamingPusher(WithMaxConcurrentOps(1))

	// Fill the semaphore
	if err := sp.AcquireSlot(context.Background()); err != nil {
		t.Fatalf("Failed to acquire slot: %v", err)
	}

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Try to acquire with cancelled context
	err := sp.AcquireSlot(ctx)
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got: %v", err)
	}

	// Cleanup
	sp.ReleaseSlot()
}

// TestStreamWithProgress verifies streaming with progress tracking
func TestStreamWithProgress(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectBytes int64
		expectCalls int
	}{
		{
			name:        "small data",
			input:       "hello world",
			expectBytes: 11,
			expectCalls: 1,
		},
		{
			name:        "large data",
			input:       strings.Repeat("a", 100*1024),
			expectBytes: 100 * 1024,
			expectCalls: 4, // 100KB / 32KB buffer = ~4 calls
		},
		{
			name:        "empty data",
			input:       "",
			expectBytes: 0,
			expectCalls: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := NewStreamingPusher()
			r := strings.NewReader(tt.input)
			w := &bytes.Buffer{}

			var progressCalls int
			var lastProgress int64

			err := sp.StreamWithProgress(
				context.Background(),
				r,
				w,
				tt.expectBytes,
				func(written int64) {
					progressCalls++
					lastProgress = written
				},
			)

			if err != nil {
				t.Fatalf("StreamWithProgress failed: %v", err)
			}

			if w.Len() != len(tt.input) {
				t.Errorf("Written bytes mismatch: got %d, want %d", w.Len(), len(tt.input))
			}

			if w.String() != tt.input {
				t.Errorf("Content mismatch")
			}

			if lastProgress != tt.expectBytes {
				t.Errorf("Last progress = %d, want %d", lastProgress, tt.expectBytes)
			}

			if tt.expectCalls > 0 && progressCalls < 1 {
				t.Errorf("Expected at least 1 progress call, got %d", progressCalls)
			}
		})
	}
}

// TestStreamWithProgressCancellation verifies cancellation during streaming
func TestStreamWithProgressCancellation(t *testing.T) {
	sp := NewStreamingPusher()

	// Create a slow reader
	slowReader := &slowReader{
		data:  []byte(strings.Repeat("a", 1024*1024)), // 1MB
		delay: 50 * time.Millisecond,
	}

	w := &bytes.Buffer{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := sp.StreamWithProgress(ctx, slowReader, w, 1024*1024, nil)

	if err != context.DeadlineExceeded {
		t.Errorf("Expected context.DeadlineExceeded, got: %v", err)
	}
}

// TestStreamWithProgressNilProgress verifies streaming works without progress callback
func TestStreamWithProgressNilProgress(t *testing.T) {
	sp := NewStreamingPusher()
	input := "test data without progress tracking"
	r := strings.NewReader(input)
	w := &bytes.Buffer{}

	err := sp.StreamWithProgress(context.Background(), r, w, int64(len(input)), nil)

	if err != nil {
		t.Fatalf("StreamWithProgress failed: %v", err)
	}

	if w.String() != input {
		t.Errorf("Content mismatch")
	}
}

// TestStreamWithProgressWriteError verifies error handling on write failures
func TestStreamWithProgressWriteError(t *testing.T) {
	sp := NewStreamingPusher()
	r := strings.NewReader("test data")
	w := &errorWriter{err: errors.New("write failed")}

	err := sp.StreamWithProgress(context.Background(), r, w, 100, nil)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err.Error() != "write failed" {
		t.Errorf("Expected 'write failed', got: %v", err)
	}
}

// TestConnectionPoolGetPut verifies connection pool operations
func TestConnectionPoolGetPut(t *testing.T) {
	cp := NewConnectionPool()

	// Test Get on empty pool
	_, exists := cp.Get("registry.example.com")
	if exists {
		t.Error("Expected connection not to exist in empty pool")
	}

	// Test Put
	conn := &PooledConnection{
		conn:     "mock-connection",
		lastUsed: time.Now().Unix(),
		useCount: 1,
		registry: "registry.example.com",
	}
	cp.Put("registry.example.com", conn)

	// Test Get after Put
	retrieved, exists := cp.Get("registry.example.com")
	if !exists {
		t.Fatal("Expected connection to exist after Put")
	}
	if retrieved.registry != "registry.example.com" {
		t.Errorf("Registry mismatch: got %s, want %s", retrieved.registry, "registry.example.com")
	}
	if retrieved.useCount != 1 {
		t.Errorf("UseCount mismatch: got %d, want 1", retrieved.useCount)
	}
}

// TestConnectionPoolMultipleRegistries verifies handling of multiple registries
func TestConnectionPoolMultipleRegistries(t *testing.T) {
	cp := NewConnectionPool()

	registries := []string{
		"registry1.example.com",
		"registry2.example.com",
		"registry3.example.com",
	}

	// Put connections for all registries
	for i, reg := range registries {
		conn := &PooledConnection{
			conn:     i,
			lastUsed: time.Now().Unix(),
			useCount: i + 1,
			registry: reg,
		}
		cp.Put(reg, conn)
	}

	// Verify all connections exist
	for i, reg := range registries {
		conn, exists := cp.Get(reg)
		if !exists {
			t.Errorf("Connection for %s not found", reg)
		}
		if conn.useCount != i+1 {
			t.Errorf("UseCount for %s: got %d, want %d", reg, conn.useCount, i+1)
		}
	}
}

// TestConnectionPoolClose verifies pool closure
func TestConnectionPoolClose(t *testing.T) {
	cp := NewConnectionPool()

	// Add some connections
	for i := 0; i < 5; i++ {
		conn := &PooledConnection{
			conn:     i,
			registry: "registry.example.com",
		}
		cp.Put("registry.example.com", conn)
	}

	// Close pool
	if err := cp.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	// Verify pool is empty
	_, exists := cp.Get("registry.example.com")
	if exists {
		t.Error("Connection still exists after Close")
	}
}

// TestConnectionPoolConcurrency verifies thread-safe operations
func TestConnectionPoolConcurrency(t *testing.T) {
	cp := NewConnectionPool()
	done := make(chan bool, 3)

	// Concurrent Put operations
	go func() {
		for i := 0; i < 1000; i++ {
			conn := &PooledConnection{
				conn:     i,
				registry: "registry1.example.com",
			}
			cp.Put("registry1.example.com", conn)
		}
		done <- true
	}()

	// Concurrent Get operations
	go func() {
		for i := 0; i < 1000; i++ {
			cp.Get("registry1.example.com")
		}
		done <- true
	}()

	// Concurrent operations on different registry
	go func() {
		for i := 0; i < 1000; i++ {
			conn := &PooledConnection{
				conn:     i,
				registry: "registry2.example.com",
			}
			cp.Put("registry2.example.com", conn)
			cp.Get("registry2.example.com")
		}
		done <- true
	}()

	// Wait for all goroutines
	for i := 0; i < 3; i++ {
		<-done
	}
}

// BenchmarkStreamingPusherGetBuffer benchmarks buffer acquisition
func BenchmarkStreamingPusherGetBuffer(b *testing.B) {
	sp := NewStreamingPusher()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := sp.GetBuffer()
		sp.PutBuffer(buf)
	}
}

// BenchmarkStreamingPusherStreamWithProgress benchmarks streaming operations
func BenchmarkStreamingPusherStreamWithProgress(b *testing.B) {
	sp := NewStreamingPusher()
	data := bytes.Repeat([]byte("a"), 1024*1024) // 1MB

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(data)
		w := io.Discard
		sp.StreamWithProgress(context.Background(), r, w, int64(len(data)), nil)
	}
}

// BenchmarkConnectionPoolGetPut benchmarks pool operations
func BenchmarkConnectionPoolGetPut(b *testing.B) {
	cp := NewConnectionPool()
	conn := &PooledConnection{
		conn:     "test",
		registry: "registry.example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cp.Put("registry.example.com", conn)
		cp.Get("registry.example.com")
	}
}

// BenchmarkStreamingPusherConcurrency benchmarks concurrent slot acquisition
func BenchmarkStreamingPusherConcurrency(b *testing.B) {
	sp := NewStreamingPusher()
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sp.AcquireSlot(ctx)
			sp.ReleaseSlot()
		}
	})
}

// slowReader simulates a slow reader for testing cancellation
type slowReader struct {
	data  []byte
	pos   int
	delay time.Duration
}

func (sr *slowReader) Read(p []byte) (n int, err error) {
	if sr.pos >= len(sr.data) {
		return 0, io.EOF
	}

	time.Sleep(sr.delay)

	n = copy(p, sr.data[sr.pos:])
	sr.pos += n
	return n, nil
}

// errorWriter simulates a writer that returns errors
type errorWriter struct {
	err error
}

func (ew *errorWriter) Write(p []byte) (n int, err error) {
	return 0, ew.err
}
