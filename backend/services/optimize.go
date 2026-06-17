package services

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"mole-wails/backend/models"
)

const optimizeTaskTimeout = 5 * time.Minute

const (
	optimizeStatusRunning = "running"
	optimizeStatusSuccess = "success"
	optimizeStatusFailed  = "failed"
)

type optimizeTaskDef struct {
	action       string
	name         string
	description  string
	requiresSudo bool
}

var optimizeTasks = []optimizeTaskDef{
	{
		action:       "system_maintenance",
		name:         "System Database Maintenance",
		description:  "Rebuild LaunchServices, refresh DNS and verify Spotlight",
		requiresSudo: true,
	},
	{
		action:       "cache_refresh",
		name:         "Finder and Safari Cache Refresh",
		description:  "Refresh QuickLook, icon services and Safari caches",
		requiresSudo: false,
	},
	{
		action:       "maintenance_scripts",
		name:         "System Log Rotation",
		description:  "Rotate and compress system logs with newsyslog",
		requiresSudo: true,
	},
	{
		action:       "radio_refresh",
		name:         "Radio Refresh",
		description:  "Refresh wireless and Bluetooth related services",
		requiresSudo: true,
	},
	{
		action:       "saved_state_cleanup",
		name:         "Saved State Cleanup",
		description:  "Remove stale application saved state data",
		requiresSudo: false,
	},
	{
		action:       "swap_cleanup",
		name:         "Virtual Memory Refresh",
		description:  "Reset swap files and dynamic pager service",
		requiresSudo: true,
	},
	{
		action:       "startup_cache",
		name:         "Startup Cache",
		description:  "Refresh startup and boot caches",
		requiresSudo: true,
	},
	{
		action:       "local_snapshots",
		name:         "Local Snapshots",
		description:  "Prune local Time Machine snapshots",
		requiresSudo: true,
	},
	{
		action:       "fix_broken_configs",
		name:         "Broken Configs",
		description:  "Repair known stale configuration states",
		requiresSudo: false,
	},
	{
		action:       "network_optimization",
		name:         "Network Stack Optimization",
		description:  "Refresh DNS, rebuild ARP and restart mDNSResponder",
		requiresSudo: true,
	},
	{
		action:       "sqlite_vacuum",
		name:         "SQLite Database Optimization",
		description:  "Optimize user databases with VACUUM to reduce size",
		requiresSudo: false,
	},
	{
		action:       "cloudshell",
		name:         "CloudShell Diagnostics",
		description:  "Report CloudShell or AliEntSafe CPU activity",
		requiresSudo: false,
	},
	{
		action:       "syspolicyd",
		name:         "Syspolicyd Diagnostics",
		description:  "Report syspolicyd verification activity",
		requiresSudo: false,
	},
	{
		action:       "windowserver",
		name:         "WindowServer Diagnostics",
		description:  "Report WindowServer resource activity",
		requiresSudo: false,
	},
	{
		action:       "spotlight",
		name:         "Spotlight Diagnostics",
		description:  "Report Spotlight indexing activity",
		requiresSudo: false,
	},
	{
		action:       "coresim_disk_images",
		name:         "CoreSimulator Images",
		description:  "Report CoreSimulator disk image usage",
		requiresSudo: false,
	},
}

var legacyOptimizeTaskIDs = map[string]string{
	"rebuild_caches":   "system_maintenance",
	"reset_network":    "network_optimization",
	"refresh_ui":       "cache_refresh",
	"clean_logs":       "maintenance_scripts",
	"restart_pager":    "swap_cleanup",
	"rebuild_services": "system_maintenance",
}

type OptimizeService struct {
	scriptsPath string
	ctx         context.Context
}

func NewOptimizeService(scriptsPath string) *OptimizeService {
	return &OptimizeService{
		scriptsPath: scriptsPath,
	}
}

func (s *OptimizeService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *OptimizeService) emitProgress(taskID, message string, percent int, status string) {
	if s.ctx == nil {
		return
	}

	runtime.EventsEmit(s.ctx, "optimize:progress", models.OptimizeProgress{
		Task:    taskID,
		Message: message,
		Percent: percent,
		Status:  status,
	})
}

func (s *OptimizeService) resolveTaskID(taskID string) (string, bool) {
	if mapped, ok := legacyOptimizeTaskIDs[taskID]; ok {
		taskID = mapped
	}

	for _, task := range optimizeTasks {
		if task.action == taskID {
			return task.action, true
		}
	}

	return "", false
}

func (s *OptimizeService) taskName(taskID string) string {
	for _, task := range optimizeTasks {
		if task.action == taskID {
			return task.name
		}
	}
	return taskID
}

// GetTasks returns available optimization tasks
func (s *OptimizeService) GetTasks() ([]models.OptimizationTask, error) {
	tasks := make([]models.OptimizationTask, 0, len(optimizeTasks))
	for _, task := range optimizeTasks {
		tasks = append(tasks, models.OptimizationTask{
			ID:           task.action,
			Name:         task.name,
			Description:  task.description,
			Enabled:      true,
			RequiresSudo: task.requiresSudo,
		})
	}

	return tasks, nil
}

func (s *OptimizeService) runTask(action string) (string, error) {
	scriptPath := filepath.Join(s.scriptsPath, "bin", "optimize_task.sh")

	ctx, cancel := context.WithTimeout(context.Background(), optimizeTaskTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/bin/bash", scriptPath, action)
	cmd.Dir = s.scriptsPath
	cmd.Stdin = nil
	cmd.Env = append(os.Environ(), "MOLE_WAILS_NONINTERACTIVE=1")

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		message := strings.TrimSpace(output.String())
		if ctx.Err() == context.DeadlineExceeded {
			if message == "" {
				message = "Task timed out"
			}
			return message, fmt.Errorf("%s timed out", action)
		}
		if message == "" {
			message = err.Error()
		}
		return message, fmt.Errorf("%s: %w", action, err)
	}

	return strings.TrimSpace(output.String()), nil
}

func lastMeaningfulLine(output string) string {
	lines := strings.Split(output, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(stripANSI(lines[i]))
		if line != "" {
			return line
		}
	}
	return ""
}

func stripANSI(line string) string {
	var b strings.Builder
	b.Grow(len(line))
	escaped := false
	for _, r := range line {
		if escaped {
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				escaped = false
			}
			continue
		}
		if r == '\x1b' {
			escaped = true
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func (s *OptimizeService) PreviewOptimizations(taskIDs []string) (models.DryRunPreview, error) {
	if len(taskIDs) == 0 {
		return models.DryRunPreview{Entries: []models.DryRunEntry{}}, nil
	}
	resolved, resolutionErrors := s.resolveTaskIDs(taskIDs)
	if len(resolutionErrors) > 0 && len(resolved) == 0 {
		return models.DryRunPreview{}, errors.New(strings.Join(resolutionErrors, "; "))
	}
	entries := make([]models.DryRunEntry, 0, len(resolved))
	for _, action := range resolved {
		entries = append(entries, models.DryRunEntry{
			Action: action,
			Detail: fmt.Sprintf("Would run %s", s.taskName(action)),
		})
	}
	return models.DryRunPreview{Entries: entries}, nil
}

func (s *OptimizeService) resolveTaskIDs(taskIDs []string) ([]string, []string) {
	resolved := make([]string, 0, len(taskIDs))
	seen := make(map[string]bool)
	var errors []string
	for _, taskID := range taskIDs {
		action, ok := s.resolveTaskID(taskID)
		if !ok {
			errors = append(errors, fmt.Sprintf("unknown optimization task: %s", taskID))
			continue
		}
		if seen[action] {
			continue
		}
		seen[action] = true
		resolved = append(resolved, action)
	}
	return resolved, errors
}

// ExecuteOptimizations runs selected optimization tasks
func (s *OptimizeService) ExecuteOptimizations(taskIDs []string, dryRun bool) error {
	if len(taskIDs) == 0 {
		return fmt.Errorf("no optimization tasks selected")
	}

	if dryRun {
		preview, err := s.PreviewOptimizations(taskIDs)
		if err != nil {
			return err
		}
		if s.ctx != nil {
			runtime.EventsEmit(s.ctx, "optimize:complete", models.OptimizeResult{TasksCompleted: len(preview.Entries)})
		}
		return nil
	}

	resolved, resolutionErrors := s.resolveTaskIDs(taskIDs)
	if len(resolved) == 0 {
		return errors.New(strings.Join(resolutionErrors, "; "))
	}

	totalTasks := len(resolved)
	completed := 0
	errors := append([]string{}, resolutionErrors...)

	s.emitProgress("", "Starting optimization...", 0, optimizeStatusRunning)

	for index, action := range resolved {
		taskName := s.taskName(action)
		startPercent := (index * 100) / totalTasks
		s.emitProgress(action, fmt.Sprintf("Running %s...", taskName), startPercent, optimizeStatusRunning)

		output, err := s.runTask(action)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", taskName, err))
			message := lastMeaningfulLine(output)
			if message == "" {
				message = fmt.Sprintf("%s failed", taskName)
			}
			s.emitProgress(action, message, startPercent, optimizeStatusFailed)
			continue
		}

		completed++
		endPercent := (completed * 100) / totalTasks
		message := lastMeaningfulLine(output)
		if message == "" {
			message = fmt.Sprintf("%s complete", taskName)
		}
		s.emitProgress(action, message, endPercent, optimizeStatusSuccess)
	}

	s.emitProgress("", "Optimization complete", 100, optimizeStatusSuccess)

	result := models.OptimizeResult{
		TasksCompleted: completed,
		Errors:         errors,
	}

	if s.ctx != nil {
		runtime.EventsEmit(s.ctx, "optimize:complete", result)
	}

	if len(errors) > 0 && completed == 0 {
		return fmt.Errorf("all optimization tasks failed")
	}

	return nil
}

// GetWhitelist returns optimization tasks in whitelist
func (s *OptimizeService) GetWhitelist() ([]string, error) {
	whitelistPath := filepath.Join(os.Getenv("HOME"), ".config", "mole", "optimize_whitelist")

	data, err := os.ReadFile(whitelistPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read whitelist: %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	var whitelist []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			whitelist = append(whitelist, line)
		}
	}

	return whitelist, nil
}

// UpdateWhitelist updates optimization whitelist
func (s *OptimizeService) UpdateWhitelist(tasks []string) error {
	whitelistPath := filepath.Join(os.Getenv("HOME"), ".config", "mole", "optimize_whitelist")

	configDir := filepath.Dir(whitelistPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	content := strings.Join(tasks, "\n")
	if err := os.WriteFile(whitelistPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write whitelist: %w", err)
	}

	return nil
}
