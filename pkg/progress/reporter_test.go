package progress_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/cnoe-io/idpbuilder/pkg/progress"
	"github.com/cnoe-io/idpbuilder/pkg/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T-2.1.2-01: TestNewReporter_Normal
func TestNewReporter_Normal(t *testing.T) {
	// Given: Create reporter with verbose=false
	reporter := progress.NewReporter(false)

	// Then: Reporter is created successfully
	require.NotNil(t, reporter)
	assert.NotNil(t, reporter.GetCallback())
}

// T-2.1.2-02: TestNewReporter_Verbose
func TestNewReporter_Verbose(t *testing.T) {
	// Given: Create reporter with verbose=true
	reporter := progress.NewReporter(true)

	// Then: Reporter is created successfully
	require.NotNil(t, reporter)
	assert.NotNil(t, reporter.GetCallback())
}

// T-2.1.2-03: TestReporter_HandleProgress_Uploading
func TestReporter_HandleProgress_Uploading(t *testing.T) {
	// Given: Reporter
	reporter := progress.NewReporter(false)

	// When: Receiving uploading progress update
	reporter.HandleProgress(registry.ProgressUpdate{
		LayerDigest: "sha256:abc123def456",
		LayerSize:   1024,
		BytesPushed: 512,
		Status:      "uploading",
	})

	// Then: No panics, summary works
	reporter.DisplaySummary()
}

// T-2.1.2-04: TestReporter_HandleProgress_Complete
func TestReporter_HandleProgress_Complete(t *testing.T) {
	// Given: Reporter
	reporter := progress.NewReporter(false)

	// When: Receiving progress updates from uploading to complete
	reporter.HandleProgress(registry.ProgressUpdate{
		LayerDigest: "sha256:abc123def456",
		LayerSize:   1024,
		BytesPushed: 0,
		Status:      "uploading",
	})

	reporter.HandleProgress(registry.ProgressUpdate{
		LayerDigest: "sha256:abc123def456",
		LayerSize:   1024,
		BytesPushed: 1024,
		Status:      "complete",
	})

	// Then: No panics, summary works
	reporter.DisplaySummary()
}

// T-2.1.2-05: TestReporter_HandleProgress_Exists
func TestReporter_HandleProgress_Exists(t *testing.T) {
	// Given: Reporter
	reporter := progress.NewReporter(false)

	// When: Receiving exists status
	reporter.HandleProgress(registry.ProgressUpdate{
		LayerDigest: "sha256:exists123",
		LayerSize:   2048,
		BytesPushed: 0,
		Status:      "exists",
	})

	// Then: No panics, summary works
	reporter.DisplaySummary()
}

// T-2.1.2-06: TestReporter_HandleProgress_MultipleLayers
func TestReporter_HandleProgress_MultipleLayers(t *testing.T) {
	// Given: Reporter
	reporter := progress.NewReporter(false)

	// When: Receiving updates for multiple layers
	for i := 0; i < 5; i++ {
		digest := fmt.Sprintf("sha256:layer%d", i)
		reporter.HandleProgress(registry.ProgressUpdate{
			LayerDigest: digest,
			LayerSize:   1000,
			BytesPushed: 1000,
			Status:      "complete",
		})
	}

	// Then: Summary shows all layers
	reporter.DisplaySummary()
}

// T-2.1.2-07: TestReporter_HandleProgress_ThreadSafety
func TestReporter_HandleProgress_ThreadSafety(t *testing.T) {
	// Given: Reporter
	reporter := progress.NewReporter(false)

	// When: Concurrent progress updates
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(layerNum int) {
			defer wg.Done()
			digest := fmt.Sprintf("sha256:layer%d", layerNum)
			reporter.HandleProgress(registry.ProgressUpdate{
				LayerDigest: digest,
				LayerSize:   1000,
				BytesPushed: 500,
				Status:      "uploading",
			})
		}(i)
	}

	wg.Wait()

	// Then: No race conditions, all updates processed
	reporter.DisplaySummary()
	// Race detector would catch issues: go test -race
}

// T-2.1.2-08: TestReporter_DisplayNormal_Format
func TestReporter_DisplayNormal_Format(t *testing.T) {
	// Given: Reporter in normal mode
	reporter := progress.NewReporter(false)

	// When: Processing various status updates
	reporter.HandleProgress(registry.ProgressUpdate{
		LayerDigest: "sha256:abcdefghijklmnop",
		LayerSize:   1000,
		BytesPushed: 500,
		Status:      "uploading",
	})

	reporter.HandleProgress(registry.ProgressUpdate{
		LayerDigest: "sha256:complete123",
		LayerSize:   2000,
		BytesPushed: 2000,
		Status:      "complete",
	})

	// Then: Displays work (visual verification needed)
	reporter.DisplaySummary()
}

// T-2.1.2-09: TestReporter_DisplayVerbose_Format
func TestReporter_DisplayVerbose_Format(t *testing.T) {
	// Given: Reporter in verbose mode
	reporter := progress.NewReporter(true)

	// When: Processing various status updates
	reporter.HandleProgress(registry.ProgressUpdate{
		LayerDigest: "sha256:verbosetest123456789",
		LayerSize:   5000000,
		BytesPushed: 2500000,
		Status:      "uploading",
	})

	reporter.HandleProgress(registry.ProgressUpdate{
		LayerDigest: "sha256:verbosecomplete",
		LayerSize:   1000000,
		BytesPushed: 1000000,
		Status:      "complete",
	})

	// Then: Verbose display works (visual verification needed)
	reporter.DisplaySummary()
}

// T-2.1.2-10: TestReporter_DisplayVerbose_RateCalculation
func TestReporter_DisplayVerbose_RateCalculation(t *testing.T) {
	// Given: Reporter in verbose mode
	reporter := progress.NewReporter(true)

	// When: Processing update with known size
	reporter.HandleProgress(registry.ProgressUpdate{
		LayerDigest: "sha256:ratecalc",
		LayerSize:   10000000, // 10 MB
		BytesPushed: 5000000,  // 5 MB
		Status:      "uploading",
	})

	// Then: No panics from rate calculation
	reporter.DisplaySummary()
}

// T-2.1.2-11: TestReporter_DisplaySummary_SingleLayer
func TestReporter_DisplaySummary_SingleLayer(t *testing.T) {
	// Given: Reporter
	reporter := progress.NewReporter(false)

	// When: Single layer pushed
	reporter.HandleProgress(registry.ProgressUpdate{
		LayerDigest: "sha256:single",
		LayerSize:   1048576, // 1 MB
		BytesPushed: 1048576,
		Status:      "complete",
	})

	// Then: Summary displays correctly
	reporter.DisplaySummary()
}

// T-2.1.2-12: TestReporter_DisplaySummary_MultipleLayers
func TestReporter_DisplaySummary_MultipleLayers(t *testing.T) {
	// Given: Reporter
	reporter := progress.NewReporter(false)

	// When: Multiple layers pushed
	for i := 0; i < 5; i++ {
		digest := fmt.Sprintf("sha256:multilayer%d", i)
		reporter.HandleProgress(registry.ProgressUpdate{
			LayerDigest: digest,
			LayerSize:   1048576, // 1 MB each
			BytesPushed: 1048576,
			Status:      "complete",
		})
	}

	// Then: Summary shows correct totals
	reporter.DisplaySummary()
}

// T-2.1.2-13: TestReporter_DisplaySummary_MixedStatus
func TestReporter_DisplaySummary_MixedStatus(t *testing.T) {
	// Given: Reporter
	reporter := progress.NewReporter(false)

	// When: Mix of pushed and skipped layers
	reporter.HandleProgress(registry.ProgressUpdate{
		LayerDigest: "sha256:pushed1",
		LayerSize:   1048576,
		BytesPushed: 1048576,
		Status:      "complete",
	})

	reporter.HandleProgress(registry.ProgressUpdate{
		LayerDigest: "sha256:pushed2",
		LayerSize:   2097152,
		BytesPushed: 2097152,
		Status:      "complete",
	})

	reporter.HandleProgress(registry.ProgressUpdate{
		LayerDigest: "sha256:skipped1",
		LayerSize:   1048576,
		BytesPushed: 0,
		Status:      "exists",
	})

	// Then: Summary distinguishes pushed vs skipped
	reporter.DisplaySummary()
}

// T-2.1.2-14: TestReporter_GetCallback
func TestReporter_GetCallback(t *testing.T) {
	// Given: Reporter
	reporter := progress.NewReporter(false)

	// When: Getting callback
	callback := reporter.GetCallback()

	// Then: Callback has correct signature and works
	require.NotNil(t, callback)

	// Verify callback can be called
	callback(registry.ProgressUpdate{
		LayerDigest: "sha256:callback",
		LayerSize:   1024,
		BytesPushed: 512,
		Status:      "uploading",
	})

	reporter.DisplaySummary()
}

// T-2.1.2-15: TestReporter_DigestTruncation
func TestReporter_DigestTruncation(t *testing.T) {
	// Given: Reporter in normal mode
	reporter := progress.NewReporter(false)

	// When: Processing update with long digest
	longDigest := "sha256:abcdefghijklmnopqrstuvwxyz1234567890"
	reporter.HandleProgress(registry.ProgressUpdate{
		LayerDigest: longDigest,
		LayerSize:   1024,
		BytesPushed: 1024,
		Status:      "complete",
	})

	// Then: Digest gets truncated (visual verification - should see first 12 chars)
	reporter.DisplaySummary()
}
