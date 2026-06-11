package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"

	"mole-wails/backend/models"
)

const (
	historyDefaultLimit      = 20
	historyMaxLimit          = 200
	historyScriptName        = "history.sh"
	historyJSONFlag          = "--json"
	historyLimitFlag         = "--limit"
	historyShellPath         = "/bin/bash"
	historyCommandErrorLimit = 1000
)

type HistoryService struct {
	scriptsPath string
	ctx         context.Context
}

type historyScriptResult struct {
	Logs      models.HistoryLogs      `json:"logs"`
	Limit     int                     `json:"limit"`
	Sessions  []historyScriptSession  `json:"sessions"`
	Deletions []historyScriptDeletion `json:"deletions"`
}

type historyScriptSession struct {
	Command        string                     `json:"command"`
	StartedAt      string                     `json:"started_at"`
	EndedAt        string                     `json:"ended_at"`
	Items          int                        `json:"items"`
	Size           string                     `json:"size"`
	OperationCount int                        `json:"operation_count"`
	Actions        models.HistoryActionCounts `json:"actions"`
}

type historyScriptDeletion struct {
	Timestamp string `json:"timestamp"`
	Mode      string `json:"mode"`
	Status    string `json:"status"`
	SizeKB    *int   `json:"size_kb"`
	Path      string `json:"path"`
}

func NewHistoryService(scriptsPath string) *HistoryService {
	return &HistoryService{scriptsPath: scriptsPath}
}

func (s *HistoryService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *HistoryService) GetHistory(limit int) (models.HistoryResult, error) {
	normalizedLimit := normalizeHistoryLimit(limit)
	scriptPath := filepath.Join(s.scriptsPath, "bin", historyScriptName)
	if err := VerifyBundledScript(scriptPath); err != nil {
		return models.HistoryResult{}, err
	}

	cmd := exec.Command(historyShellPath, scriptPath, historyJSONFlag, historyLimitFlag, strconv.Itoa(normalizedLimit))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return models.HistoryResult{}, fmt.Errorf("failed to load history: %s", trimCommandOutput(output))
	}

	var scriptResult historyScriptResult
	if err := json.Unmarshal(output, &scriptResult); err != nil {
		return models.HistoryResult{}, fmt.Errorf("failed to parse history: %w", err)
	}

	return mapHistoryScriptResult(scriptResult, normalizedLimit), nil
}

func normalizeHistoryLimit(limit int) int {
	if limit <= 0 {
		return historyDefaultLimit
	}
	if limit > historyMaxLimit {
		return historyMaxLimit
	}
	return limit
}

func mapHistoryScriptResult(scriptResult historyScriptResult, fallbackLimit int) models.HistoryResult {
	result := models.HistoryResult{
		Logs:      scriptResult.Logs,
		Limit:     scriptResult.Limit,
		Sessions:  make([]models.HistorySession, 0, len(scriptResult.Sessions)),
		Deletions: make([]models.HistoryDeletion, 0, len(scriptResult.Deletions)),
	}
	if result.Limit == 0 {
		result.Limit = fallbackLimit
	}

	for _, session := range scriptResult.Sessions {
		result.Sessions = append(result.Sessions, models.HistorySession{
			Command:        session.Command,
			StartedAt:      session.StartedAt,
			EndedAt:        session.EndedAt,
			Items:          session.Items,
			Size:           session.Size,
			OperationCount: session.OperationCount,
			Actions:        session.Actions,
		})
	}

	for _, deletion := range scriptResult.Deletions {
		result.Deletions = append(result.Deletions, models.HistoryDeletion{
			Timestamp: deletion.Timestamp,
			Mode:      deletion.Mode,
			Status:    deletion.Status,
			SizeKB:    deletion.SizeKB,
			Path:      deletion.Path,
		})
	}

	return result
}

func trimCommandOutput(output []byte) string {
	text := string(output)
	if len(text) > historyCommandErrorLimit {
		return text[:historyCommandErrorLimit]
	}
	return text
}
