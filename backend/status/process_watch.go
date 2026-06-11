package status

import (
	"fmt"
	"sync"
	"time"

	"mole-wails/backend/models"
)

const (
	defaultProcessWatchThreshold = 100
	defaultProcessWatchWindow    = 5 * time.Minute
)

type ProcessWatcher struct {
	mu        sync.Mutex
	threshold float64
	window    time.Duration
	enabled   bool
	tracked   map[processKey]processTrack
}

type processKey struct {
	pid     int
	ppid    int
	command string
}

type processTrack struct {
	startedAt time.Time
	alerting  bool
	process   ProcessInfo
}

func NewProcessWatcher() *ProcessWatcher {
	return &ProcessWatcher{
		threshold: defaultProcessWatchThreshold,
		window:    defaultProcessWatchWindow,
		enabled:   true,
		tracked:   make(map[processKey]processTrack),
	}
}

func (w *ProcessWatcher) Config() models.ProcessWatchConfig {
	w.mu.Lock()
	defer w.mu.Unlock()
	return models.ProcessWatchConfig{
		Threshold:     w.threshold,
		WindowSeconds: int(w.window.Seconds()),
		Enabled:       w.enabled,
	}
}

func (w *ProcessWatcher) SetConfig(config models.ProcessWatchConfig) error {
	if config.Threshold < 0 {
		return fmt.Errorf("invalid CPU threshold: %.2f", config.Threshold)
	}
	if config.WindowSeconds <= 0 {
		return fmt.Errorf("invalid process watch window: %d", config.WindowSeconds)
	}

	w.mu.Lock()
	defer w.mu.Unlock()
	w.threshold = config.Threshold
	w.window = time.Duration(config.WindowSeconds) * time.Second
	w.enabled = config.Enabled
	if !w.enabled {
		w.tracked = make(map[processKey]processTrack)
	}
	return nil
}

func (w *ProcessWatcher) Update(now time.Time, processes []ProcessInfo) []models.ProcessAlert {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.enabled {
		return []models.ProcessAlert{}
	}

	seen := make(map[processKey]bool, len(processes))
	for _, process := range processes {
		if process.PID <= 0 {
			continue
		}
		key := processKey{pid: process.PID, ppid: process.PPID, command: process.Command}
		seen[key] = true
		track, exists := w.tracked[key]
		if process.CPU < w.threshold {
			delete(w.tracked, key)
			continue
		}
		if !exists {
			track = processTrack{startedAt: now}
		}
		track.process = process
		track.alerting = now.Sub(track.startedAt) >= w.window
		w.tracked[key] = track
	}

	for key := range w.tracked {
		if !seen[key] {
			delete(w.tracked, key)
		}
	}

	alerts := make([]models.ProcessAlert, 0)
	for _, track := range w.tracked {
		if !track.alerting {
			continue
		}
		alerts = append(alerts, models.ProcessAlert{
			PID:           track.process.PID,
			Name:          track.process.Name,
			Command:       track.process.Command,
			CPUPercent:    track.process.CPU,
			Threshold:     w.threshold,
			WindowSeconds: int(w.window.Seconds()),
			TriggeredAt:   track.startedAt.Add(w.window),
			Status:        "active",
		})
	}
	return alerts
}
