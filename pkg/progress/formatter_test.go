package progress

import (
	"testing"
	"time"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{100, "100 B"},
		{1023, "1023 B"},
		{1024, "1.0 KiB"},
		{1536, "1.5 KiB"},
		{1048576, "1.0 MiB"},
		{1073741824, "1.0 GiB"},
		{1099511627776, "1.0 TiB"},
	}

	for _, test := range tests {
		result := formatBytes(test.bytes)
		if result != test.expected {
			t.Errorf("formatBytes(%d) = %s, expected %s", test.bytes, result, test.expected)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{100 * time.Millisecond, "100ms"},
		{500 * time.Millisecond, "500ms"},
		{999 * time.Millisecond, "999ms"},
		{1 * time.Second, "1.0s"},
		{1500 * time.Millisecond, "1.5s"},
		{30 * time.Second, "30.0s"},
		{59 * time.Second, "59.0s"},
		{1 * time.Minute, "1.0m"},
		{90 * time.Second, "1.5m"},
		{5 * time.Minute, "5.0m"},
	}

	for _, test := range tests {
		result := formatDuration(test.duration)
		if result != test.expected {
			t.Errorf("formatDuration(%v) = %s, expected %s", test.duration, result, test.expected)
		}
	}
}

func TestFormatEdgeCases(t *testing.T) {
	// Test zero duration
	result := formatDuration(0)
	if result != "0ms" {
		t.Errorf("formatDuration(0) = %s, expected 0ms", result)
	}

	// Test very large bytes
	largeBytes := int64(1024 * 1024 * 1024 * 1024 * 1024) // 1 PiB
	result = formatBytes(largeBytes)
	if result != "1.0 PiB" {
		t.Errorf("formatBytes(1PiB) = %s, expected 1.0 PiB", result)
	}
}