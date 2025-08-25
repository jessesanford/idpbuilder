package progress

import "time"

type ProgressStatus string

const (
	StatusPending  ProgressStatus = "pending"
	StatusRunning  ProgressStatus = "running"
	StatusComplete ProgressStatus = "complete"
	StatusError    ProgressStatus = "error"
)

type ProgressEvent struct {
	ID        string                 `json:"id"`
	Operation string                 `json:"operation"`
	Status    ProgressStatus         `json:"status"`
	Message   string                 `json:"message"`
	Timestamp time.Time              `json:"timestamp"`
	Progress  int                    `json:"progress,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

type ProgressReporter interface {
	ReportProgress(event ProgressEvent) error
	ReportError(id, operation string, err error) error
	ReportComplete(id, operation, message string) error
}

type BuildProgress struct {
	BuildID       string `json:"build_id"`
	Stage         string `json:"stage"`
	StepTotal     int    `json:"step_total,omitempty"`
	StepCurrent   int    `json:"step_current,omitempty"`
	BytesTotal    int64  `json:"bytes_total,omitempty"`
	BytesCurrent  int64  `json:"bytes_current,omitempty"`
	LayersTotal   int    `json:"layers_total,omitempty"`
	LayersCurrent int    `json:"layers_current,omitempty"`
}

type PushProgress struct {
	Repository   string `json:"repository"`
	Tag          string `json:"tag"`
	LayerID      string `json:"layer_id,omitempty"`
	BytesTotal   int64  `json:"bytes_total,omitempty"`
	BytesCurrent int64  `json:"bytes_current,omitempty"`
}
