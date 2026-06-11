package services

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"mole-wails/backend/models"
)

type TouchIDService struct {
	scriptsPath string
	ctx         context.Context
}

func NewTouchIDService(scriptsPath string) *TouchIDService {
	return &TouchIDService{
		scriptsPath: scriptsPath,
	}
}

func (s *TouchIDService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// GetStatus checks if Touch ID for sudo is enabled
func (s *TouchIDService) GetStatus() (*models.TouchIDStatus, error) {
	configPath := "/etc/pam.d/sudo"
	pamModulePath := "/usr/lib/pam/pam_tid.so.2"

	// Check if pam_tid module exists
	_, err := os.Stat(pamModulePath)
	available := err == nil

	// Check if Touch ID is enabled in config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read sudoers file: %w", err)
	}

	content := string(data)
	enabled := strings.Contains(content, "pam_tid.so")

	status := &models.TouchIDStatus{
		Enabled:       enabled,
		Available:     available,
		Status:        "Disabled",
		PamModulePath: pamModulePath,
		ConfigPath:    configPath,
	}

	if enabled {
		status.Status = "Enabled"
	}

	return status, nil
}

// Enable enables Touch ID for sudo
func (s *TouchIDService) Enable(dryRun bool) error {
	if dryRun {
		_, err := s.Preview("enable")
		return err
	}
	scriptPath := filepath.Join(s.scriptsPath, "bin", "touchid.sh")
	if err := VerifyBundledScript(scriptPath); err != nil {
		return err
	}

	cmd := exec.Command("/bin/bash", scriptPath, "enable")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to enable Touch ID: %w, output: %s", err, string(output))
	}

	return nil
}

// Disable disables Touch ID for sudo
func (s *TouchIDService) Disable(dryRun bool) error {
	if dryRun {
		_, err := s.Preview("disable")
		return err
	}
	scriptPath := filepath.Join(s.scriptsPath, "bin", "touchid.sh")
	if err := VerifyBundledScript(scriptPath); err != nil {
		return err
	}

	cmd := exec.Command("/bin/bash", scriptPath, "disable")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to disable Touch ID: %w, output: %s", err, string(output))
	}

	return nil
}

func (s *TouchIDService) Preview(action string) (models.DryRunPreview, error) {
	if action != "enable" && action != "disable" {
		return models.DryRunPreview{}, fmt.Errorf("unknown Touch ID action: %s", action)
	}
	scriptPath := filepath.Join(s.scriptsPath, "bin", "touchid.sh")
	if err := VerifyBundledScript(scriptPath); err != nil {
		return models.DryRunPreview{}, err
	}
	cmd := exec.Command("/bin/bash", scriptPath, action, "--dry-run")
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return models.DryRunPreview{}, fmt.Errorf("Touch ID preview failed: %w, output: %s", err, output.String())
	}
	entries := []models.DryRunEntry{}
	for _, line := range strings.Split(output.String(), "\n") {
		line = strings.TrimSpace(stripANSI(line))
		if line == "" {
			continue
		}
		entries = append(entries, models.DryRunEntry{Action: action, Detail: line})
	}
	return models.DryRunPreview{Entries: entries}, nil
}
