package services

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"mole-wails/backend/models"
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
		action:       "swap_cleanup",
		name:         "Virtual Memory Refresh",
		description:  "Reset swap files and dynamic pager service",
		requiresSudo: true,
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
}

var legacyOptimizeTaskIDs = map[string]string{
	"rebuild_caches":    "system_maintenance",
	"reset_network":     "network_optimization",
	"refresh_ui":        "cache_refresh",
	"clean_logs":        "maintenance_scripts",
	"restart_pager":     "swap_cleanup",
	"rebuild_services":  "system_maintenance",
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

func (s *OptimizeService) emitProgress(taskID, message string, percent int) {
	if s.ctx == nil {
		return
	}

	runtime.EventsEmit(s.ctx, "optimize:progress", models.OptimizeProgress{
		Task:    taskID,
		Message: message,
		Percent: percent,
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

	cmd := exec.Command("/bin/bash", scriptPath, action)
	cmd.Dir = s.scriptsPath
	cmd.Stdin = nil

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		message := strings.TrimSpace(output.String())
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

// ExecuteOptimizations runs selected optimization tasks
func (s *OptimizeService) ExecuteOptimizations(taskIDs []string) error {
	if len(taskIDs) == 0 {
		return fmt.Errorf("no optimization tasks selected")
	}

	resolved := make([]string, 0, len(taskIDs))
	seen := make(map[string]bool)
	for _, taskID := range taskIDs {
		action, ok := s.resolveTaskID(taskID)
		if !ok {
			return fmt.Errorf("unknown optimization task: %s", taskID)
		}
		if seen[action] {
			continue
		}
		seen[action] = true
		resolved = append(resolved, action)
	}

	totalTasks := len(resolved)
	completed := 0
	var errors []string

	s.emitProgress("", "Starting optimization...", 0)

	for index, action := range resolved {
		taskName := s.taskName(action)
		startPercent := (index * 100) / totalTasks
		s.emitProgress(action, fmt.Sprintf("Running %s...", taskName), startPercent)

		output, err := s.runTask(action)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", taskName, err))
			message := lastMeaningfulLine(output)
			if message == "" {
				message = fmt.Sprintf("%s failed", taskName)
			}
			s.emitProgress(action, message, startPercent)
			continue
		}

		completed++
		endPercent := (completed * 100) / totalTasks
		message := lastMeaningfulLine(output)
		if message == "" {
			message = fmt.Sprintf("%s complete", taskName)
		}
		s.emitProgress(action, message, endPercent)
	}

	s.emitProgress("", "Optimization complete", 100)

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
