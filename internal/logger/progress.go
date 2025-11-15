package logger

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

type ProgressBar struct {
	total     int
	current   int
	width     int
	mu        sync.Mutex
	writer    io.Writer
	startTime time.Time
	lastPrint time.Time
}

func NewProgressBar(total int, writer io.Writer) *ProgressBar {
	return &ProgressBar{
		total:     total,
		current:   0,
		width:     50,
		writer:    writer,
		startTime: time.Now(),
		lastPrint: time.Now(),
	}
}

func (pb *ProgressBar) Increment(description string) {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	pb.current++
	pb.render(description)
}

func (pb *ProgressBar) SetCurrent(current int, description string) {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	pb.current = current
	pb.render(description)
}

func (pb *ProgressBar) Finish(description string) {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	pb.current = pb.total
	pb.render(description)
	fmt.Fprintln(pb.writer)
}

func (pb *ProgressBar) render(description string) {

	now := time.Now()
	if now.Sub(pb.lastPrint) < 100*time.Millisecond && pb.current < pb.total {
		return
	}
	pb.lastPrint = now

	percentage := float64(pb.current) / float64(pb.total)
	filled := int(percentage * float64(pb.width))

	bar := strings.Repeat("█", filled)
	empty := strings.Repeat("░", pb.width-filled)

	elapsed := time.Since(pb.startTime)
	var eta string
	if pb.current > 0 {
		avgTime := elapsed / time.Duration(pb.current)
		remaining := time.Duration(pb.total-pb.current) * avgTime
		eta = fmt.Sprintf("ETA: %s", formatDuration(remaining))
	} else {
		eta = "ETA: calculating..."
	}

	fmt.Fprintf(pb.writer, "\r[%s%s] %3.0f%% (%d/%d) %s %s",
		bar, empty, percentage*100, pb.current, pb.total, description, eta)
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)

	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h > 0 {
		return fmt.Sprintf("%dh%dm%ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm%ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}
