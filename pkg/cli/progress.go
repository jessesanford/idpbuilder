package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// ProgressReporter provides progress updates for long-running operations
type ProgressReporter interface {
	UpdateMessage(message string)
	UpdateProgress(current, total int64)
	Finish()
}

// ProgressBar implements a simple progress bar with spinner support
type ProgressBar struct {
	mu           sync.Mutex
	writer       io.Writer
	message      string
	current      int64
	total        int64
	startTime    time.Time
	spinnerChars []string
	spinnerIndex int
	finished     bool
}

// NewProgressBar creates a new progress reporter
func NewProgressBar(initialMessage string) ProgressReporter {
	pb := &ProgressBar{
		writer:       os.Stdout,
		message:      initialMessage,
		startTime:    time.Now(),
		spinnerChars: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		spinnerIndex: 0,
	}
	pb.render()
	return pb
}

// NewQuietProgress creates a progress reporter that only outputs final messages
func NewQuietProgress(initialMessage string) ProgressReporter {
	return &QuietProgress{message: initialMessage}
}

// UpdateMessage updates the current progress message
func (pb *ProgressBar) UpdateMessage(message string) {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	
	if pb.finished {
		return
	}
	
	pb.message = message
	pb.render()
}

// UpdateProgress updates the current progress values
func (pb *ProgressBar) UpdateProgress(current, total int64) {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	
	if pb.finished {
		return
	}
	
	pb.current = current
	pb.total = total
	pb.render()
}

// Finish completes the progress reporting
func (pb *ProgressBar) Finish() {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	
	if pb.finished {
		return
	}
	
	pb.finished = true
	elapsed := time.Since(pb.startTime)
	
	// Clear line and show final message
	pb.clearLine()
	fmt.Fprintf(pb.writer, "✓ %s (took %v)\n", pb.message, elapsed.Round(time.Millisecond))
}

// render displays the current progress state
func (pb *ProgressBar) render() {
	// Clear current line
	pb.clearLine()
	
	// Show spinner
	spinner := pb.spinnerChars[pb.spinnerIndex%len(pb.spinnerChars)]
	pb.spinnerIndex++
	
	// Build progress line
	line := fmt.Sprintf("%s %s", spinner, pb.message)
	
	// Add percentage if we have total
	if pb.total > 0 {
		percent := int((pb.current * 100) / pb.total)
		line += fmt.Sprintf(" [%3d%%]", percent)
		
		// Add progress bar
		barWidth := 20
		filled := int((pb.current * int64(barWidth)) / pb.total)
		bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
		line += fmt.Sprintf(" %s", bar)
		
		// Add size info if meaningful
		if pb.total > 1024 {
			line += fmt.Sprintf(" (%s/%s)", formatBytes(pb.current), formatBytes(pb.total))
		}
	}
	
	// Add elapsed time
	elapsed := time.Since(pb.startTime)
	line += fmt.Sprintf(" (%v)", elapsed.Round(time.Second))
	
	fmt.Fprint(pb.writer, line)
}

// clearLine clears the current terminal line
func (pb *ProgressBar) clearLine() {
	fmt.Fprint(pb.writer, "\r\033[K")
}

// formatBytes formats byte counts in human-readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	units := []string{"KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

// QuietProgress implements a quiet progress reporter that only shows final results
type QuietProgress struct {
	message string
}

// UpdateMessage updates the message for quiet progress
func (qp *QuietProgress) UpdateMessage(message string) {
	qp.message = message
}

// UpdateProgress is a no-op for quiet progress
func (qp *QuietProgress) UpdateProgress(current, total int64) {
	// No-op for quiet mode
}

// Finish outputs the final message for quiet progress
func (qp *QuietProgress) Finish() {
	fmt.Printf("✓ %s\n", qp.message)
}

// MultiProgress manages multiple concurrent progress bars
type MultiProgress struct {
	mu        sync.RWMutex
	writer    io.Writer
	bars      map[string]*ProgressBar
	finished  map[string]bool
	order     []string
}

// NewMultiProgress creates a new multi-progress manager
func NewMultiProgress() *MultiProgress {
	return &MultiProgress{
		writer:   os.Stdout,
		bars:     make(map[string]*ProgressBar),
		finished: make(map[string]bool),
	}
}

// AddBar adds a new progress bar with the given name
func (mp *MultiProgress) AddBar(name, message string) ProgressReporter {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	
	bar := &ProgressBar{
		writer:       mp.writer,
		message:      message,
		startTime:    time.Now(),
		spinnerChars: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	}
	
	mp.bars[name] = bar
	mp.finished[name] = false
	mp.order = append(mp.order, name)
	
	return &MultiProgressBar{mp: mp, name: name}
}

// MultiProgressBar wraps a progress bar for multi-progress management
type MultiProgressBar struct {
	mp   *MultiProgress
	name string
}

// UpdateMessage updates the message for this progress bar
func (mpb *MultiProgressBar) UpdateMessage(message string) {
	mpb.mp.mu.Lock()
	defer mpb.mp.mu.Unlock()
	
	if bar, ok := mpb.mp.bars[mpb.name]; ok && !mpb.mp.finished[mpb.name] {
		bar.message = message
		mpb.mp.renderAll()
	}
}

// UpdateProgress updates the progress for this progress bar
func (mpb *MultiProgressBar) UpdateProgress(current, total int64) {
	mpb.mp.mu.Lock()
	defer mpb.mp.mu.Unlock()
	
	if bar, ok := mpb.mp.bars[mpb.name]; ok && !mpb.mp.finished[mpb.name] {
		bar.current = current
		bar.total = total
		mpb.mp.renderAll()
	}
}

// Finish marks this progress bar as finished
func (mpb *MultiProgressBar) Finish() {
	mpb.mp.mu.Lock()
	defer mpb.mp.mu.Unlock()
	
	mpb.mp.finished[mpb.name] = true
	mpb.mp.renderAll()
}

// renderAll renders all active progress bars
func (mp *MultiProgress) renderAll() {
	// Move cursor up to start of progress section
	activeCount := 0
	for _, finished := range mp.finished {
		if !finished {
			activeCount++
		}
	}
	
	if activeCount > 1 {
		fmt.Fprintf(mp.writer, "\033[%dA", activeCount-1)
	}
	
	// Render each bar
	for _, name := range mp.order {
		if bar, ok := mp.bars[name]; ok {
			if !mp.finished[name] {
				bar.render()
				fmt.Fprintln(mp.writer)
			}
		}
	}
}