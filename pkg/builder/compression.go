package builder

// CompressionType represents the type of compression to use.
type CompressionType int
const (
	// NoCompression indicates no compression
	NoCompression CompressionType = iota
	// GzipCompression indicates gzip compression
	GzipCompression
	// Bzip2Compression indicates bzip2 compression
	Bzip2Compression
	// XzCompression indicates xz compression
	XzCompression
)
// String returns the string representation of the compression type.
func (ct CompressionType) String() string {
	switch ct {
	case NoCompression:
		return "none"
	case GzipCompression:
		return "gzip"
	case Bzip2Compression:
		return "bzip2"
	case XzCompression:
		return "xz"
	default:
		return "unknown"
	}
}