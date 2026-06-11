package models

import "time"

// Clean service types

type CleanCategory struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Enabled        bool   `json:"enabled"`
	EstimatedMB    int64  `json:"estimatedMB"`
	EstimatedBytes int64  `json:"estimatedBytes"`
	RequiresSudo   bool   `json:"requiresSudo"`
}

type CleanProgress struct {
	Category    string `json:"category"`
	Message     string `json:"message"`
	Percent     int    `json:"percent"`
	CurrentFile string `json:"currentFile"`
	TotalFiles  int    `json:"totalFiles"`
	FilesClean  int    `json:"filesClean"`
}

type CleanResult struct {
	SpaceFreed   int64    `json:"spaceFreed"`
	FilesRemoved int      `json:"filesRemoved"`
	Categories   []string `json:"categories"`
	Errors       []string `json:"errors"`
}

// Uninstall service types

type Application struct {
	Name         string    `json:"name"`
	BundleID     string    `json:"bundleId"`
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
	Age          string    `json:"age"`
	Icon         string    `json:"icon,omitempty"`
	BrewCask     string    `json:"brewCask,omitempty"`
}

type UninstallProgress struct {
	App          string `json:"app"`
	Message      string `json:"message"`
	Percent      int    `json:"percent"`
	FilesRemoved int    `json:"filesRemoved"`
	TotalFiles   int    `json:"totalFiles"`
	SpaceFreed   int64  `json:"spaceFreed"`
}

type UninstallResult struct {
	AppsRemoved  int      `json:"appsRemoved"`
	FilesRemoved int      `json:"filesRemoved"`
	SpaceFreed   int64    `json:"spaceFreed"`
	Errors       []string `json:"errors"`
}

// Optimize service types

type OptimizationTask struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Enabled      bool   `json:"enabled"`
	RequiresSudo bool   `json:"requiresSudo"`
}

type OptimizeProgress struct {
	Task    string `json:"task"`
	Message string `json:"message"`
	Percent int    `json:"percent"`
}

type OptimizeResult struct {
	TasksCompleted int      `json:"tasksCompleted"`
	Errors         []string `json:"errors"`
}

type DryRunEntry struct {
	Action string `json:"action"`
	Path   string `json:"path,omitempty"`
	Detail string `json:"detail"`
}

type DryRunPreview struct {
	Entries []DryRunEntry `json:"entries"`
}

type PurgeArtifact struct {
	Path     string `json:"path"`
	Type     string `json:"type"`
	Size     int64  `json:"size"`
	AgeDays  int    `json:"ageDays"`
	Selected bool   `json:"selected"`
}

type PurgeProgress struct {
	CurrentPath  string `json:"currentPath"`
	ItemsScanned int    `json:"itemsScanned"`
}

type PurgeScanResult struct {
	Artifacts       []PurgeArtifact `json:"artifacts"`
	Errors          []string        `json:"errors"`
	TotalSize       int64           `json:"totalSize"`
	ArtifactCount   int             `json:"artifactCount"`
	ConfiguredPaths []string        `json:"configuredPaths"`
	MissingPaths    []string        `json:"missingPaths"`
}

type PurgeResult struct {
	SpaceFreed   int64    `json:"spaceFreed"`
	RemovedCount int      `json:"removedCount"`
	Errors       []string `json:"errors"`
}

type InstallerSource string

const (
	InstallerSourceDownloads InstallerSource = "Downloads"
	InstallerSourceDesktop   InstallerSource = "Desktop"
	InstallerSourceDocuments InstallerSource = "Documents"
	InstallerSourcePublic    InstallerSource = "Public"
	InstallerSourceLibrary   InstallerSource = "Library"
	InstallerSourceShared    InstallerSource = "Shared"
	InstallerSourceHomebrew  InstallerSource = "Homebrew"
	InstallerSourceICloud    InstallerSource = "iCloud"
	InstallerSourceMail      InstallerSource = "Mail"
	InstallerSourceTelegram  InstallerSource = "Telegram"
)

type InstallerFile struct {
	Path         string          `json:"path"`
	Size         int64           `json:"size"`
	LastModified time.Time       `json:"lastModified"`
	Source       InstallerSource `json:"source"`
	Selected     bool            `json:"selected"`
}

type InstallerProgress struct {
	CurrentPath  string `json:"currentPath"`
	ItemsScanned int    `json:"itemsScanned"`
}

type InstallerScanResult struct {
	Files     []InstallerFile `json:"files"`
	Errors    []string        `json:"errors"`
	TotalSize int64           `json:"totalSize"`
	FileCount int             `json:"fileCount"`
}

type InstallerResult struct {
	SpaceFreed   int64    `json:"spaceFreed"`
	RemovedCount int      `json:"removedCount"`
	Errors       []string `json:"errors"`
}

type HistoryActionCounts struct {
	Removed int `json:"removed"`
	Trashed int `json:"trashed"`
	Skipped int `json:"skipped"`
	Failed  int `json:"failed"`
	Rebuilt int `json:"rebuilt"`
	Other   int `json:"other"`
}

type HistorySession struct {
	Command        string              `json:"command"`
	StartedAt      string              `json:"startedAt"`
	EndedAt        string              `json:"endedAt"`
	Items          int                 `json:"items"`
	Size           string              `json:"size"`
	OperationCount int                 `json:"operationCount"`
	Actions        HistoryActionCounts `json:"actions"`
}

type HistoryDeletion struct {
	Timestamp string `json:"timestamp"`
	Mode      string `json:"mode"`
	Status    string `json:"status"`
	SizeKB    *int   `json:"sizeKb"`
	Path      string `json:"path"`
}

type HistoryLogs struct {
	Operations string `json:"operations"`
	Deletions  string `json:"deletions"`
}

type HistoryResult struct {
	Logs      HistoryLogs       `json:"logs"`
	Limit     int               `json:"limit"`
	Sessions  []HistorySession  `json:"sessions"`
	Deletions []HistoryDeletion `json:"deletions"`
}

// Analyze service types

type FileEntry struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Size int64  `json:"size"`
}

type DirEntry struct {
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	Size       int64     `json:"size"`
	IsDir      bool      `json:"isDir"`
	LastAccess time.Time `json:"lastAccess"`
	Percent    float64   `json:"percent"`
}

type ScanResult struct {
	Entries    []DirEntry  `json:"entries"`
	LargeFiles []FileEntry `json:"largeFiles"`
	TotalSize  int64       `json:"totalSize"`
	TotalItems int         `json:"totalItems"`
	Path       string      `json:"path"`
}

type ScanProgress struct {
	Path         string `json:"path"`
	ItemsScanned int    `json:"itemsScanned"`
	TotalSize    int64  `json:"totalSize"`
}

type ExternalVolume struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Available bool   `json:"available"`
}

// Status service types

type MetricsSnapshot struct {
	// Hardware info
	Hardware HardwareInfo `json:"hardware"`

	// Health score (0-100)
	Health int `json:"health"`

	// CPU metrics
	CPU CPUMetrics `json:"cpu"`

	// GPU metrics
	GPU GPUMetrics `json:"gpu"`

	// Memory metrics
	Memory MemoryMetrics `json:"memory"`

	// Disk metrics
	Disk DiskMetrics `json:"disk"`

	// Network metrics
	Network NetworkMetrics `json:"network"`

	// Battery/Power metrics
	Battery BatteryMetrics `json:"battery"`

	// Top processes
	Processes []ProcessInfo `json:"processes"`

	// Persistent high-CPU process alerts
	ProcessAlerts []ProcessAlert `json:"processAlerts"`

	// Timestamp
	Timestamp time.Time `json:"timestamp"`
}

type HardwareInfo struct {
	Model     string `json:"model"`
	Processor string `json:"processor"`
	Memory    string `json:"memory"`
	OS        string `json:"os"`
	OSVersion string `json:"osVersion"`
	Uptime    string `json:"uptime"`
}

type CPUMetrics struct {
	TotalPercent float64   `json:"totalPercent"`
	LoadAvg      []float64 `json:"loadAvg"` // 1m, 5m, 15m
	Cores        int       `json:"cores"`
	PerCore      []float64 `json:"perCore"`
	Temperature  float64   `json:"temperature"`
}

type GPUMetrics struct {
	Usage       float64 `json:"usage"`
	Temperature float64 `json:"temperature"`
	Name        string  `json:"name"`
}

type MemoryMetrics struct {
	Used      int64   `json:"used"`
	Total     int64   `json:"total"`
	Free      int64   `json:"free"`
	Available int64   `json:"available"`
	Percent   float64 `json:"percent"`
}

type DiskMetrics struct {
	Used       int64   `json:"used"`
	Total      int64   `json:"total"`
	Free       int64   `json:"free"`
	Percent    float64 `json:"percent"`
	ReadBytes  int64   `json:"readBytes"`
	WriteBytes int64   `json:"writeBytes"`
	ReadSpeed  float64 `json:"readSpeed"`  // MB/s
	WriteSpeed float64 `json:"writeSpeed"` // MB/s
}

type NetworkMetrics struct {
	Download    float64 `json:"download"` // MB/s
	Upload      float64 `json:"upload"`   // MB/s
	ProxyHost   string  `json:"proxyHost"`
	ProxyPort   string  `json:"proxyPort"`
	ProxyType   string  `json:"proxyType"`
	BluetoothOn bool    `json:"bluetoothOn"`
}

type BatteryMetrics struct {
	Level       int     `json:"level"`  // Percentage
	Status      string  `json:"status"` // Charging, Charged, Discharging
	Health      string  `json:"health"` // Normal, Replace Soon, Replace Now
	Cycles      int     `json:"cycles"`
	Temperature float64 `json:"temperature"`
	FanSpeed    int     `json:"fanSpeed"` // RPM
}

type ProcessInfo struct {
	Name       string  `json:"name"`
	PID        int     `json:"pid"`
	PPID       int     `json:"ppid"`
	Command    string  `json:"command"`
	CPUPercent float64 `json:"cpuPercent"`
	MemoryMB   int64   `json:"memoryMB"`
}

type ProcessAlert struct {
	PID           int       `json:"pid"`
	Name          string    `json:"name"`
	Command       string    `json:"command"`
	CPUPercent    float64   `json:"cpuPercent"`
	Threshold     float64   `json:"threshold"`
	WindowSeconds int       `json:"windowSeconds"`
	TriggeredAt   time.Time `json:"triggeredAt"`
	Status        string    `json:"status"`
}

type ProcessWatchConfig struct {
	Threshold     float64 `json:"threshold"`
	WindowSeconds int     `json:"windowSeconds"`
	Enabled       bool    `json:"enabled"`
}

// TouchID service types

type TouchIDStatus struct {
	Enabled       bool   `json:"enabled"`
	Available     bool   `json:"available"`
	Status        string `json:"status"`
	PamModulePath string `json:"pamModulePath"`
	ConfigPath    string `json:"configPath"`
}

// Common types

type ErrorResponse struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
