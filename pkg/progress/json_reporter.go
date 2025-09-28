package progress

import (
	"encoding/json"
	"io"
	"os"
	"sync"
	"time"
)

type jsonReporter struct {
	writer    io.Writer
	encoder   *json.Encoder
	operation *Operation
	mutex     sync.Mutex
}

// NewJSONReporter creates a progress reporter with JSON output
func NewJSONReporter(writer io.Writer) ProgressReporter {
	if writer == nil {
		writer = os.Stdout
	}
	return &jsonReporter{
		writer:  writer,
		encoder: json.NewEncoder(writer),
	}
}

func (j *jsonReporter) Start(message string) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	j.operation = &Operation{
		Name:      message,
		StartTime: time.Now(),
	}

	j.encoder.Encode(map[string]interface{}{
		"event":     "start",
		"message":   message,
		"timestamp": j.operation.StartTime,
	})
}

func (j *jsonReporter) ReportProgress(current, total int64, message string) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	if j.operation != nil {
		j.operation.Current = current
		j.operation.Total = total
	}

	data := map[string]interface{}{
		"event":   "progress",
		"current": current,
		"total":   total,
		"message": message,
	}

	if total > 0 {
		percentage := float64(current) / float64(total) * 100
		data["percentage"] = percentage
	}

	j.encoder.Encode(data)
}

func (j *jsonReporter) Complete(message string) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	var duration time.Duration
	if j.operation != nil {
		duration = time.Since(j.operation.StartTime)
	}

	j.encoder.Encode(map[string]interface{}{
		"event":    "complete",
		"message":  message,
		"duration": duration.Seconds(),
	})
}

func (j *jsonReporter) Error(err error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	if err != nil {
		j.encoder.Encode(map[string]interface{}{
			"event":   "error",
			"message": err.Error(),
		})
	}
}