package progress

import (
	"sync"
)

type multiReporter struct {
	reporters []ProgressReporter
	mutex     sync.Mutex
}

// NewMultiReporter creates a composite reporter that forwards to multiple reporters
func NewMultiReporter(reporters ...ProgressReporter) ProgressReporter {
	// Filter out nil reporters for safety
	validReporters := make([]ProgressReporter, 0, len(reporters))
	for _, r := range reporters {
		if r != nil {
			validReporters = append(validReporters, r)
		}
	}

	return &multiReporter{
		reporters: validReporters,
	}
}

func (m *multiReporter) Start(message string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, r := range m.reporters {
		r.Start(message)
	}
}

func (m *multiReporter) ReportProgress(current, total int64, message string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, r := range m.reporters {
		r.ReportProgress(current, total, message)
	}
}

func (m *multiReporter) Complete(message string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, r := range m.reporters {
		r.Complete(message)
	}
}

func (m *multiReporter) Error(err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, r := range m.reporters {
		r.Error(err)
	}
}