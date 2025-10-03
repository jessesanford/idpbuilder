// Package push provides performance optimizations for OCI image push operations
package push

import (
	"context"
	"io"
	"sync"
)

const (
	// DefaultBufferSize is the default buffer size for streaming operations
	DefaultBufferSize = 32 * 1024 // 32KB

	// DefaultChunkSize is the default chunk size for layer uploads
	DefaultChunkSize = 1024 * 1024 // 1MB

	// DefaultMaxConcurrentLayers is the default maximum number of concurrent layer uploads
	DefaultMaxConcurrentLayers = 4
)

// StreamingPusher provides optimized image streaming with minimal memory footprint
type StreamingPusher struct {
	bufferPool       *sync.Pool
	chunkSize        int
	maxConcurrentOps int
	semaphore        chan struct{}
}

// NewStreamingPusher creates a new StreamingPusher with optimized settings
func NewStreamingPusher(opts ...StreamingOption) *StreamingPusher {
	sp := &StreamingPusher{
		chunkSize:        DefaultChunkSize,
		maxConcurrentOps: DefaultMaxConcurrentLayers,
		bufferPool: &sync.Pool{
			New: func() interface{} {
				buf := make([]byte, DefaultBufferSize)
				return &buf
			},
		},
	}

	for _, opt := range opts {
		opt(sp)
	}

	sp.semaphore = make(chan struct{}, sp.maxConcurrentOps)
	return sp
}

// StreamingOption is a functional option for StreamingPusher configuration
type StreamingOption func(*StreamingPusher)

// WithChunkSize sets the chunk size for uploads
func WithChunkSize(size int) StreamingOption {
	return func(sp *StreamingPusher) {
		sp.chunkSize = size
	}
}

// WithMaxConcurrentOps sets the maximum number of concurrent operations
func WithMaxConcurrentOps(max int) StreamingOption {
	return func(sp *StreamingPusher) {
		sp.maxConcurrentOps = max
	}
}

// GetBuffer returns a buffer from the pool
func (sp *StreamingPusher) GetBuffer() *[]byte {
	return sp.bufferPool.Get().(*[]byte)
}

// PutBuffer returns a buffer to the pool
func (sp *StreamingPusher) PutBuffer(buf *[]byte) {
	sp.bufferPool.Put(buf)
}

// AcquireSlot acquires a slot for concurrent operations
func (sp *StreamingPusher) AcquireSlot(ctx context.Context) error {
	select {
	case sp.semaphore <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// ReleaseSlot releases a slot for concurrent operations
func (sp *StreamingPusher) ReleaseSlot() {
	<-sp.semaphore
}

// StreamWithProgress streams data with progress tracking
func (sp *StreamingPusher) StreamWithProgress(ctx context.Context, r io.Reader, w io.Writer, total int64, progress func(written int64)) error {
	buf := sp.GetBuffer()
	defer sp.PutBuffer(buf)

	var written int64
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, err := r.Read(*buf)
		if n > 0 {
			nw, werr := w.Write((*buf)[:n])
			if nw > 0 {
				written += int64(nw)
				if progress != nil {
					progress(written)
				}
			}
			if werr != nil {
				return werr
			}
			if nw != n {
				return io.ErrShortWrite
			}
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

// ConnectionPool manages reusable HTTP connections for registry operations
type ConnectionPool struct {
	mu    sync.RWMutex
	conns map[string]*PooledConnection
}

// PooledConnection represents a pooled connection with metadata
type PooledConnection struct {
	conn      interface{}
	lastUsed  int64
	useCount  int
	registry  string
}

// NewConnectionPool creates a new connection pool
func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{
		conns: make(map[string]*PooledConnection),
	}
}

// Get retrieves or creates a connection for the given registry
func (cp *ConnectionPool) Get(registry string) (*PooledConnection, bool) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	conn, exists := cp.conns[registry]
	return conn, exists
}

// Put stores a connection in the pool
func (cp *ConnectionPool) Put(registry string, conn *PooledConnection) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	cp.conns[registry] = conn
}

// Close closes all connections in the pool
func (cp *ConnectionPool) Close() error {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	for k := range cp.conns {
		delete(cp.conns, k)
	}
	return nil
}
