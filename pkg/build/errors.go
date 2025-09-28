package build

import (
	"errors"
	"fmt"
)

// Error constants
var (
	ErrInvalidConfig  = errors.New("invalid build configuration")
	ErrBuildFailed    = errors.New("build operation failed")
	ErrLayerAddFailed = errors.New("failed to add layer")
	ErrFinalizeFailed = errors.New("failed to finalize build")
	ErrStorageInit    = errors.New("failed to initialize storage")
)

// BuildError wraps build-related errors with context
type BuildError struct {
	Op      string // Operation that failed
	Err     error  // Underlying error
	Context string // Additional context
}

func (e *BuildError) Error() string {
	if e.Context != "" {
		return fmt.Sprintf("build error in %s: %v (context: %s)", e.Op, e.Err, e.Context)
	}
	return fmt.Sprintf("build error in %s: %v", e.Op, e.Err)
}

func (e *BuildError) Unwrap() error {
	return e.Err
}

// WrapBuildError wraps an error with build context
func WrapBuildError(op string, err error, context string) error {
	if err == nil {
		return nil
	}
	return &BuildError{
		Op:      op,
		Err:     err,
		Context: context,
	}
}

// IsBuildError checks if an error is a BuildError
func IsBuildError(err error) bool {
	var buildErr *BuildError
	return errors.As(err, &buildErr)
}