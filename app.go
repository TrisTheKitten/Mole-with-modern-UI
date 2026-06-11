package main

import (
	"context"
	"os"
	"path/filepath"
	"runtime"

	"mole-wails/backend/analyze"
	"mole-wails/backend/models"
	"mole-wails/backend/services"
	"mole-wails/backend/status"
)

const (
	scriptsDirName             = "scripts"
	scriptsHistoryRelativePath = "bin/history.sh"
	macOSResourcesScriptsPath  = "../Resources/scripts"
	installedScriptsPath       = "/Applications/Mole.app/Contents/Resources/scripts"
)

// App struct
type App struct {
	ctx context.Context

	// Services
	Clean     *services.CleanService
	Uninstall *services.UninstallService
	Optimize  *services.OptimizeService
	Purge     *services.PurgeService
	Installer *services.InstallerService
	History   *services.HistoryService
	Analyze   *analyze.Service
	Status    *status.Service
	TouchID   *services.TouchIDService
}

// NewApp creates a new App application struct
func NewApp() *App {
	// Determine scripts path
	scriptsPath := getScriptsPath()

	return &App{
		Clean:     services.NewCleanService(scriptsPath),
		Uninstall: services.NewUninstallService(scriptsPath),
		Optimize:  services.NewOptimizeService(scriptsPath),
		Purge:     services.NewPurgeService(scriptsPath),
		Installer: services.NewInstallerService(scriptsPath),
		History:   services.NewHistoryService(scriptsPath),
		Analyze:   analyze.NewService(),
		Status:    status.NewService(),
		TouchID:   services.NewTouchIDService(scriptsPath),
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Set context for all services
	a.Clean.SetContext(ctx)
	a.Uninstall.SetContext(ctx)
	a.Optimize.SetContext(ctx)
	a.Purge.SetContext(ctx)
	a.Installer.SetContext(ctx)
	a.History.SetContext(ctx)
	a.Analyze.SetContext(ctx)
	a.Status.SetContext(ctx)
	a.TouchID.SetContext(ctx)
}

// shutdown is called when the app shuts down
func (a *App) shutdown(ctx context.Context) {
	// Cleanup
	a.Status.StopMonitoring()
}

// Helper function to determine scripts path
func getScriptsPath() string {
	for _, candidate := range scriptsPathCandidates() {
		if isUsableScriptsPath(candidate) {
			return candidate
		}
	}

	return installedScriptsPath
}

func scriptsPathCandidates() []string {
	candidates := make([]string, 0, 4)

	if workingDirScriptsPath, err := filepath.Abs(scriptsDirName); err == nil {
		candidates = append(candidates, workingDirScriptsPath)
	}

	if _, sourceFile, _, ok := runtime.Caller(0); ok {
		candidates = append(candidates, filepath.Join(filepath.Dir(sourceFile), scriptsDirName))
	}

	if executablePath, err := os.Executable(); err == nil {
		executableDir := filepath.Dir(executablePath)
		candidates = append(candidates, filepath.Clean(filepath.Join(executableDir, macOSResourcesScriptsPath)))
	}

	candidates = append(candidates, installedScriptsPath)
	return candidates
}

func isUsableScriptsPath(scriptsPath string) bool {
	info, err := os.Stat(filepath.Join(scriptsPath, scriptsHistoryRelativePath))
	return err == nil && !info.IsDir()
}

// ===========================
// Clean Service Methods
// ===========================

func (a *App) CleanScanTargets() ([]models.CleanCategory, error) {
	return a.Clean.ScanTargets()
}

func (a *App) CleanExecute(categories []string, dryRun bool) error {
	return a.Clean.ExecuteClean(categories, dryRun)
}

func (a *App) CleanGetWhitelist() ([]string, error) {
	return a.Clean.GetWhitelist()
}

func (a *App) CleanUpdateWhitelist(paths []string) error {
	return a.Clean.UpdateWhitelist(paths)
}

// ===========================
// Uninstall Service Methods
// ===========================

func (a *App) UninstallScanApps(forceRescan bool) ([]models.Application, error) {
	return a.Uninstall.ScanApplications(forceRescan)
}

func (a *App) UninstallApps(bundleIDs []string) error {
	return a.Uninstall.UninstallApps(bundleIDs, false)
}

func (a *App) UninstallAppsWithDryRun(bundleIDs []string, dryRun bool) error {
	return a.Uninstall.UninstallApps(bundleIDs, dryRun)
}

func (a *App) UninstallGetRelatedFiles(bundleID string) ([]string, error) {
	return a.Uninstall.GetRelatedFiles(bundleID)
}

func (a *App) UninstallPreview(bundleIDs []string) (models.DryRunPreview, error) {
	return a.Uninstall.PreviewUninstall(bundleIDs)
}

// ===========================
// Optimize Service Methods
// ===========================

func (a *App) OptimizeGetTasks() ([]models.OptimizationTask, error) {
	return a.Optimize.GetTasks()
}

func (a *App) OptimizeExecute(taskIDs []string) error {
	return a.Optimize.ExecuteOptimizations(taskIDs, false)
}

func (a *App) OptimizeExecuteWithDryRun(taskIDs []string, dryRun bool) error {
	return a.Optimize.ExecuteOptimizations(taskIDs, dryRun)
}

func (a *App) OptimizePreview(taskIDs []string) (models.DryRunPreview, error) {
	return a.Optimize.PreviewOptimizations(taskIDs)
}

func (a *App) OptimizeGetWhitelist() ([]string, error) {
	return a.Optimize.GetWhitelist()
}

func (a *App) OptimizeUpdateWhitelist(tasks []string) error {
	return a.Optimize.UpdateWhitelist(tasks)
}

// ===========================
// Analyze Service Methods
// ===========================

func (a *App) AnalyzeScanDirectory(path string) (*models.ScanResult, error) {
	return a.Analyze.ScanDirectory(path)
}

func (a *App) AnalyzeGetLargeFiles(path string, limit int) ([]models.FileEntry, error) {
	return a.Analyze.GetLargeFiles(path, limit)
}

func (a *App) AnalyzeDeletePath(path string) error {
	return a.Analyze.DeletePath(path)
}

func (a *App) AnalyzeOpenInFinder(path string) error {
	return a.Analyze.OpenInFinder(path)
}

func (a *App) AnalyzePickDirectory(defaultPath string) (string, error) {
	return a.Analyze.PickDirectory(defaultPath)
}

func (a *App) AnalyzeListExternalVolumes() ([]models.ExternalVolume, error) {
	return a.Analyze.ListExternalVolumes()
}

func (a *App) AnalyzeScanExternalVolume(path string) (*models.ScanResult, error) {
	return a.Analyze.ScanExternalVolume(path)
}

// ===========================
// Status Service Methods
// ===========================

func (a *App) StatusGetMetrics() (*models.MetricsSnapshot, error) {
	return a.Status.GetMetrics()
}

func (a *App) StatusStartMonitoring(interval int) error {
	return a.Status.StartMonitoring(interval)
}

func (a *App) StatusStopMonitoring() {
	a.Status.StopMonitoring()
}

func (a *App) StatusGetProcessWatchConfig() models.ProcessWatchConfig {
	return a.Status.GetProcessWatchConfig()
}

func (a *App) StatusSetProcessWatchConfig(config models.ProcessWatchConfig) error {
	return a.Status.SetProcessWatchConfig(config)
}

// ===========================
// TouchID Service Methods
// ===========================

func (a *App) TouchIDGetStatus() (*models.TouchIDStatus, error) {
	return a.TouchID.GetStatus()
}

func (a *App) TouchIDEnable() error {
	return a.TouchID.Enable(false)
}

func (a *App) TouchIDEnableWithDryRun(dryRun bool) error {
	return a.TouchID.Enable(dryRun)
}

func (a *App) TouchIDDisable() error {
	return a.TouchID.Disable(false)
}

func (a *App) TouchIDDisableWithDryRun(dryRun bool) error {
	return a.TouchID.Disable(dryRun)
}

func (a *App) TouchIDPreview(action string) (models.DryRunPreview, error) {
	return a.TouchID.Preview(action)
}

// ===========================
// Purge Service Methods
// ===========================

func (a *App) PurgeScan() (models.PurgeScanResult, error) {
	return a.Purge.ScanTargets()
}

func (a *App) PurgeExecute(paths []string) (models.PurgeResult, error) {
	return a.Purge.ExecutePurge(paths)
}

func (a *App) PurgeGetPaths() ([]string, error) {
	return a.Purge.GetPaths()
}

func (a *App) PurgeUpdatePaths(paths []string) error {
	return a.Purge.UpdatePaths(paths)
}

// ===========================
// Installer Service Methods
// ===========================

func (a *App) InstallerScan() (models.InstallerScanResult, error) {
	return a.Installer.ScanInstallers()
}

func (a *App) InstallerRemove(paths []string) (models.InstallerResult, error) {
	return a.Installer.RemoveInstallers(paths)
}

// ===========================
// History Service Methods
// ===========================

func (a *App) HistoryGet(limit int) (models.HistoryResult, error) {
	return a.History.GetHistory(limit)
}
