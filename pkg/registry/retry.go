package registry

import (
	"fmt"
	"log"
	"time"
)

// retryWithExponentialBackoff executes the given operation with exponential backoff retry logic
func retryWithExponentialBackoff(operation func() error, operationName, target string) error {
	const (
		maxRetries    = 3
		initialDelay  = time.Second
		backoffFactor = 2.0
	)

	var lastErr error
	delay := initialDelay

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			log.Printf("Retrying %s for %s (attempt %d/%d) after %v delay",
				operationName, target, attempt+1, maxRetries, delay)
			time.Sleep(delay)
			delay = time.Duration(float64(delay) * backoffFactor)
		}

		lastErr = operation()
		if lastErr == nil {
			if attempt > 0 {
				log.Printf("Successfully completed %s for %s after %d retries",
					operationName, target, attempt)
			}
			return nil
		}

		log.Printf("Failed %s for %s (attempt %d/%d): %v",
			operationName, target, attempt+1, maxRetries, lastErr)
	}

	return fmt.Errorf("failed %s for %s after %d attempts: %v",
		operationName, target, maxRetries, lastErr)
}
<<<<<<< HEAD
=======

func NewWithRetry(registry Registry, policy *RetryPolicy) *WithRetry {
	if policy == nil {
		policy = DefaultRetryPolicy()
	}
	return &WithRetry{registry, policy}
}

func (r *WithRetry) Push(ctx context.Context, image string, content io.Reader) error {
	return RetryWithPolicy(ctx, r.policy, func() error { return r.registry.Push(ctx, image, content) })
}

func (r *WithRetry) List(ctx context.Context) ([]string, error) {
	var result []string
	err := RetryWithPolicy(ctx, r.policy, func() error { 
		var e error
		result, e = r.registry.List(ctx)
		return e
	})
	return result, err
}

func (r *WithRetry) Exists(ctx context.Context, repository string) (bool, error) {
	var result bool
	err := RetryWithPolicy(ctx, r.policy, func() error { 
		var e error
		result, e = r.registry.Exists(ctx, repository)
		return e
	})
	return result, err
}

func (r *WithRetry) Delete(ctx context.Context, repository string) error {
	return RetryWithPolicy(ctx, r.policy, func() error { return r.registry.Delete(ctx, repository) })
}

func (r *WithRetry) Close() error { return r.registry.Close() }
// retryWithExponentialBackoff is a wrapper for backward compatibility
// This function is used by split-001 code
func retryWithExponentialBackoff(operation func() error, operationName string, details string) error {
	ctx := context.Background()
	policy := DefaultRetryPolicy()
	return RetryWithPolicy(ctx, policy, operation)
}
>>>>>>> origin/idpbuilder-oci-build-push/phase2/wave2/integration-20250916-002118
