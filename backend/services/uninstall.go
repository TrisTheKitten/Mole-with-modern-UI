package services

import (
	"bufio"
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

const unknownBundleID = "unknown"

const (
	stepLocate  = "locate"
	stepRemove  = "remove"
	stepCleanup = "cleanup"

	stepStatusRunning = "running"
	stepStatusSuccess = "success"
	stepStatusFailed  = "failed"

	uninstallStepsPerApp = 3
)

// uninstallStepReporter surfaces a single uninstall step transition to the
// frontend. It is supplied by the orchestration loop so uninstallApplication
// stays focused on the file work while emission/percent logic lives in one place.
type uninstallStepReporter func(step, status, message string)

type UninstallService struct {
	scriptsPath string
	ctx         context.Context
}

func NewUninstallService(scriptsPath string) *UninstallService {
	return &UninstallService{
		scriptsPath: scriptsPath,
	}
}

func (s *UninstallService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// ScanApplications scans installed applications
func (s *UninstallService) ScanApplications(forceRescan bool) ([]models.Application, error) {
	var apps []models.Application

	appDirs := []string{
		"/Applications",
		filepath.Join(os.Getenv("HOME"), "Applications"),
	}

	for _, dir := range appDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !strings.HasSuffix(entry.Name(), ".app") {
				continue
			}

			appPath := filepath.Join(dir, entry.Name())
			info, err := entry.Info()
			if err != nil {
				continue
			}

			// Get app size
			size, _ := s.getDirSize(appPath)

			// Calculate age
			age := s.calculateAge(info.ModTime())

			apps = append(apps, models.Application{
				Name:         strings.TrimSuffix(entry.Name(), ".app"),
				BundleID:     s.getBundleID(appPath),
				Path:         appPath,
				Size:         size,
				LastModified: info.ModTime(),
				Age:          age,
				BrewCask:     s.detectBrewCask(appPath),
			})
		}
	}

	return apps, nil
}

// UninstallApps uninstalls selected applications
func (s *UninstallService) UninstallApps(appIdentifiers []string, dryRun bool) error {
	if dryRun {
		preview, err := s.PreviewUninstall(appIdentifiers)
		if err != nil {
			return err
		}
		if s.ctx != nil {
			runtime.EventsEmit(s.ctx, "uninstall:complete", models.UninstallResult{
				AppsRemoved:  0,
				FilesRemoved: len(preview.Entries),
			})
		}
		return nil
	}

	apps, err := s.resolveApplications(appIdentifiers)
	if err != nil {
		return err
	}

	result := models.UninstallResult{}
	totalSteps := len(apps) * uninstallStepsPerApp
	completedSteps := 0

	for _, app := range apps {
		report := func(step, status, message string) {
			if status != stepStatusRunning {
				completedSteps++
			}
			percent := 0
			if totalSteps > 0 {
				percent = (completedSteps * 100) / totalSteps
			}
			s.emitStep(app.Name, step, status, message, percent, result.FilesRemoved, result.SpaceFreed)
		}

		removedFiles, freedBytes, err := s.uninstallApplication(app, report)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %s", app.Name, err.Error()))
			continue
		}

		result.AppsRemoved++
		result.FilesRemoved += removedFiles
		result.SpaceFreed += freedBytes
	}

	if s.ctx != nil {
		s.emitProgress("", "Finished", 100, result.FilesRemoved, result.SpaceFreed)
		runtime.EventsEmit(s.ctx, "uninstall:complete", result)
	}

	return nil
}

func (s *UninstallService) PreviewUninstall(appIdentifiers []string) (models.DryRunPreview, error) {
	apps, err := s.resolveApplications(appIdentifiers)
	if err != nil {
		return models.DryRunPreview{}, err
	}

	entries := make([]models.DryRunEntry, 0)
	for _, app := range apps {
		entries = append(entries, models.DryRunEntry{Action: "remove", Path: app.Path, Detail: app.Name})
		files, err := s.GetRelatedFiles(app.Path)
		if err != nil {
			return models.DryRunPreview{}, fmt.Errorf("preview failed for %s: %w", app.Name, err)
		}
		for _, file := range files {
			entries = append(entries, models.DryRunEntry{Action: "remove", Path: file, Detail: file})
		}
	}
	return models.DryRunPreview{Entries: entries}, nil
}

// GetRelatedFiles finds all files related to an application
func (s *UninstallService) GetRelatedFiles(identifier string) ([]string, error) {
	bundleID := identifier
	appName := ""
	if app, ok := s.resolveApplication(identifier); ok {
		bundleID = app.BundleID
		appName = app.Name
	}

	if bundleID == "" || bundleID == unknownBundleID {
		return []string{}, nil
	}

	searchPaths := []string{
		filepath.Join(os.Getenv("HOME"), "Library", "Application Support"),
		filepath.Join(os.Getenv("HOME"), "Library", "Caches"),
		filepath.Join(os.Getenv("HOME"), "Library", "Preferences"),
		filepath.Join(os.Getenv("HOME"), "Library", "Logs"),
		filepath.Join(os.Getenv("HOME"), "Library", "Cookies"),
	}

	var relatedFiles []string

	for _, searchPath := range searchPaths {
		entries, err := os.ReadDir(searchPath)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if matchesApplicationFile(entry.Name(), bundleID, appName) {
				relatedFiles = append(relatedFiles, filepath.Join(searchPath, entry.Name()))
			}
		}
	}

	return relatedFiles, nil
}

// Helper functions

func (s *UninstallService) resolveApplications(identifiers []string) ([]models.Application, error) {
	if len(identifiers) == 0 {
		return nil, errors.New("select at least one app")
	}

	allApps, err := s.ScanApplications(false)
	if err != nil {
		return nil, err
	}

	apps := make([]models.Application, 0, len(identifiers))
	missing := make([]string, 0)

	for _, identifier := range identifiers {
		app, ok := findApplication(allApps, identifier)
		if !ok {
			missing = append(missing, identifier)
			continue
		}
		apps = append(apps, app)
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("applications not found: %s", strings.Join(missing, ", "))
	}

	return apps, nil
}

func (s *UninstallService) resolveApplication(identifier string) (models.Application, bool) {
	allApps, err := s.ScanApplications(false)
	if err != nil {
		return models.Application{}, false
	}
	return findApplication(allApps, identifier)
}

func findApplication(apps []models.Application, identifier string) (models.Application, bool) {
	for _, app := range apps {
		if app.Path == identifier {
			return app, true
		}
	}
	for _, app := range apps {
		if app.BundleID != "" && app.BundleID != unknownBundleID && app.BundleID == identifier {
			return app, true
		}
	}
	return models.Application{}, false
}

func (s *UninstallService) uninstallApplication(app models.Application, report uninstallStepReporter) (int, int64, error) {
	report(stepLocate, stepStatusRunning, "Scanning application files")

	if app.Path == "" || !strings.HasSuffix(app.Path, ".app") {
		report(stepLocate, stepStatusFailed, "This doesn't look like a valid app")
		return 0, 0, errors.New("invalid app path")
	}
	if _, err := os.Stat(app.Path); err != nil {
		report(stepLocate, stepStatusFailed, "Already removed — nothing left here")
		return 0, 0, fmt.Errorf("app is not available: %w", err)
	}

	relatedFiles, err := s.GetRelatedFiles(app.Path)
	if err != nil {
		report(stepLocate, stepStatusFailed, "Couldn't read its support files")
		return 0, 0, err
	}
	report(stepLocate, stepStatusSuccess, locatedSummary(len(relatedFiles)))

	spaceFreed := app.Size
	for _, relatedFile := range relatedFiles {
		size, err := s.getPathSize(relatedFile)
		if err == nil {
			spaceFreed += size
		}
	}

	if app.BrewCask != "" {
		report(stepRemove, stepStatusRunning, "Removing with Homebrew")
		if err := s.uninstallBrewCask(app); err != nil {
			report(stepRemove, stepStatusFailed, friendlyRemovalError(err))
			return 0, 0, err
		}
	} else {
		report(stepRemove, stepStatusRunning, "Deleting application files")
		if err := os.RemoveAll(app.Path); err != nil {
			report(stepRemove, stepStatusFailed, "Couldn't delete it — quit the app and try again")
			return 0, 0, err
		}
	}

	if _, err := os.Stat(app.Path); err == nil {
		report(stepRemove, stepStatusFailed, "Some files are locked — admin access may be needed")
		return 0, 0, errors.New("app is still installed")
	} else if !os.IsNotExist(err) {
		report(stepRemove, stepStatusFailed, "Couldn't confirm it was removed")
		return 0, 0, err
	}
	report(stepRemove, stepStatusSuccess, "Application deleted")

	report(stepCleanup, stepStatusRunning, "Clearing caches & preferences")
	leftoverRemoved := 0
	for _, relatedFile := range relatedFiles {
		if !isUserLibraryPath(relatedFile) {
			continue
		}
		if err := os.RemoveAll(relatedFile); err == nil {
			leftoverRemoved++
		}
	}
	report(stepCleanup, stepStatusSuccess, cleanupSummary(leftoverRemoved))

	return leftoverRemoved + 1, spaceFreed, nil
}

func locatedSummary(count int) string {
	switch count {
	case 0:
		return "No extra files found"
	case 1:
		return "Found 1 related item"
	default:
		return fmt.Sprintf("Found %d related items", count)
	}
}

func cleanupSummary(count int) string {
	switch count {
	case 0:
		return "Nothing left to clean"
	case 1:
		return "Removed 1 leftover file"
	default:
		return fmt.Sprintf("Removed %d leftover files", count)
	}
}

func friendlyRemovalError(err error) string {
	const maxLen = 120
	message := strings.TrimSpace(err.Error())
	if idx := strings.IndexByte(message, '\n'); idx != -1 {
		message = strings.TrimSpace(message[:idx])
	}
	if message == "" {
		return "Homebrew couldn't remove it"
	}
	if len(message) > maxLen {
		message = strings.TrimSpace(message[:maxLen-1]) + "…"
	}
	return message
}

func (s *UninstallService) uninstallBrewCask(app models.Application) error {
	if _, err := exec.LookPath("brew"); err != nil {
		return errors.New("Homebrew is not available")
	}

	cmd := exec.Command("brew", "uninstall", "--cask", "--zap", app.BrewCask)
	output, err := cmd.CombinedOutput()
	if err != nil {
		message := strings.TrimSpace(string(output))
		if message == "" {
			message = err.Error()
		}
		return errors.New(message)
	}

	return nil
}

func (s *UninstallService) emitProgress(app string, message string, percent int, filesRemoved int, spaceFreed int64) {
	s.emitStep(app, "", "", message, percent, filesRemoved, spaceFreed)
}

func (s *UninstallService) emitStep(app, step, status, message string, percent, filesRemoved int, spaceFreed int64) {
	if s.ctx == nil {
		return
	}

	runtime.EventsEmit(s.ctx, "uninstall:progress", models.UninstallProgress{
		App:          app,
		Step:         step,
		Status:       status,
		Message:      message,
		Percent:      percent,
		FilesRemoved: filesRemoved,
		SpaceFreed:   spaceFreed,
	})
}

func matchesApplicationFile(name string, bundleID string, appName string) bool {
	lowerName := strings.ToLower(name)
	if strings.Contains(lowerName, strings.ToLower(bundleID)) {
		return true
	}
	if appName == "" {
		return false
	}
	normalizedAppName := strings.ToLower(strings.ReplaceAll(appName, " ", ""))
	normalizedName := strings.ToLower(strings.ReplaceAll(name, " ", ""))
	return normalizedAppName != "" && strings.Contains(normalizedName, normalizedAppName)
}

func isUserLibraryPath(path string) bool {
	homeLibrary := filepath.Join(os.Getenv("HOME"), "Library")
	cleanPath := filepath.Clean(path)
	return cleanPath == homeLibrary || strings.HasPrefix(cleanPath, homeLibrary+string(os.PathSeparator))
}

func (s *UninstallService) getPathSize(path string) (int64, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return 0, err
	}
	if !info.IsDir() {
		return info.Size(), nil
	}
	return s.getDirSize(path)
}

func (s *UninstallService) getDirSize(path string) (int64, error) {
	var size int64

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}

func (s *UninstallService) getBundleID(appPath string) string {
	plistPath := filepath.Join(appPath, "Contents", "Info.plist")

	cmd := exec.Command("plutil", "-extract", "CFBundleIdentifier", "raw", plistPath)
	output, err := cmd.Output()
	if err != nil {
		return unknownBundleID
	}

	return strings.TrimSpace(string(output))
}

func (s *UninstallService) detectBrewCask(appPath string) string {
	resolvedPath, err := filepath.EvalSymlinks(appPath)
	if err != nil {
		resolvedPath = appPath
	}
	parts := strings.Split(resolvedPath, string(os.PathSeparator))
	for i, part := range parts {
		if part == "Caskroom" && i+1 < len(parts) {
			return parts[i+1]
		}
	}

	cmd := exec.Command("brew", "list", "--cask")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	appName := strings.TrimSuffix(filepath.Base(appPath), ".app")
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		cask := strings.TrimSpace(scanner.Text())
		if cask == "" {
			continue
		}
		if strings.EqualFold(cask, appName) || strings.EqualFold(strings.ReplaceAll(cask, "-", " "), appName) {
			return cask
		}
	}
	return ""
}

func (s *UninstallService) calculateAge(modTime time.Time) string {
	duration := time.Since(modTime)
	days := int(duration.Hours() / 24)

	if days < 7 {
		return "Recent"
	} else if days < 30 {
		return "< 1 month"
	} else if days < 90 {
		return "< 3 months"
	} else if days < 180 {
		return "< 6 months"
	} else {
		return "Old"
	}
}
