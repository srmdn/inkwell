package rebuild

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type Status string

const (
	StatusIdle    Status = "idle"
	StatusRunning Status = "running"
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
)

// State holds the current rebuild state. Safe for concurrent access.
type State struct {
	mu         sync.RWMutex
	Status     Status
	Output     string
	StartedAt  time.Time
	FinishedAt time.Time
	Error      string
}

func (s *State) snapshot() StatusResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return StatusResponse{
		Status:     s.Status,
		Output:     s.Output,
		StartedAt:  s.StartedAt,
		FinishedAt: s.FinishedAt,
		Error:      s.Error,
	}
}

// StatusResponse is the JSON-serialisable view of State.
type StatusResponse struct {
	Status     Status    `json:"status"`
	Output     string    `json:"output,omitempty"`
	StartedAt  time.Time `json:"started_at,omitempty"`
	FinishedAt time.Time `json:"finished_at,omitempty"`
	Error      string    `json:"error,omitempty"`
}

// Rebuilder manages async theme builds.
type Rebuilder struct {
	state        *State
	themeDir     string
	buildCmd     string
	themeService string // optional systemd service to restart after build
}

func New(themeDir, buildCmd, themeService string) *Rebuilder {
	return &Rebuilder{
		state:        &State{Status: StatusIdle},
		themeDir:     themeDir,
		buildCmd:     buildCmd,
		themeService: themeService,
	}
}

// Trigger starts an async rebuild. Returns false if one is already running.
func (rb *Rebuilder) Trigger() bool {
	rb.state.mu.Lock()
	if rb.state.Status == StatusRunning {
		rb.state.mu.Unlock()
		return false
	}
	rb.state.Status = StatusRunning
	rb.state.Output = ""
	rb.state.Error = ""
	rb.state.StartedAt = time.Now()
	rb.state.FinishedAt = time.Time{}
	rb.state.mu.Unlock()

	go rb.run()
	return true
}

// GetStatus returns a snapshot of the current rebuild state.
func (rb *Rebuilder) GetStatus() StatusResponse {
	return rb.state.snapshot()
}

func (rb *Rebuilder) run() {
	output, err := rb.build()

	rb.state.mu.Lock()
	defer rb.state.mu.Unlock()

	rb.state.Output = output
	rb.state.FinishedAt = time.Now()

	if err != nil {
		rb.state.Status = StatusFailed
		rb.state.Error = err.Error()
		return
	}

	if rb.themeService != "" {
		if err := restartService(rb.themeService); err != nil {
			rb.state.Status = StatusFailed
			rb.state.Error = fmt.Sprintf("build succeeded but service restart failed: %v", err)
			return
		}
	}

	rb.state.Status = StatusSuccess
}

func (rb *Rebuilder) build() (string, error) {
	parts := strings.Fields(rb.buildCmd)
	if len(parts) == 0 {
		return "", fmt.Errorf("THEME_BUILD_CMD is empty")
	}

	if _, err := os.Stat(rb.themeDir); err != nil {
		return "", fmt.Errorf("theme directory %q not found — set THEME_DIR in your .env", rb.themeDir)
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = rb.themeDir
	cmd.Env = append(os.Environ(), "ASTRO_TELEMETRY_DISABLED=1")

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if err := cmd.Run(); err != nil {
		return buf.String(), fmt.Errorf("build command failed: %w", err)
	}
	return buf.String(), nil
}

func restartService(name string) error {
	cmd := exec.Command("sudo", "systemctl", "restart", name)
	var buf bytes.Buffer
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, buf.String())
	}
	return nil
}
