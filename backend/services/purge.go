package services

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mole-wails/backend/models"
)

const purgeMinimumAgeDays = 7

var purgeDefaultPaths = []string{"~/Projects", "~/GitHub", "~/dev"}

var purgeArtifactNames = map[string]string{
	"node_modules": "Node modules",
	"target":       "Rust target",
	".build":       "Swift build",
	"build":        "Build output",
	"dist":         "Distribution",
	"venv":         "Python venv",
}

type PurgeService struct {
	scriptsPath string
	ctx         context.Context
}

func NewPurgeService(scriptsPath string) *PurgeService {
	return &PurgeService{scriptsPath: scriptsPath}
}

func (s *PurgeService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *PurgeService) ScanTargets() (models.PurgeScanResult, error) {
	paths, err := s.GetPaths()
	if err != nil {
		return models.PurgeScanResult{}, err
	}

	emitter := NewRuntimeOperationEmitter(s.ctx, "purge")
	result := models.PurgeScanResult{ConfiguredPaths: paths}
	scanned := 0

	for _, configuredPath := range paths {
		root := expandUserPath(configuredPath)
		info, err := os.Stat(root)
		if err != nil || !info.IsDir() {
			result.MissingPaths = append(result.MissingPaths, configuredPath)
			result.Errors = append(result.Errors, fmt.Sprintf("%s: unavailable", configuredPath))
			continue
		}

		walkErr := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", path, walkErr))
				return filepath.SkipDir
			}
			if !entry.IsDir() {
				return nil
			}

			scanned++
			emitter.Progress(models.PurgeProgress{CurrentPath: path, ItemsScanned: scanned})

			artifactType, ok := purgeArtifactNames[entry.Name()]
			if !ok {
				return nil
			}
			info, err := entry.Info()
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", path, err))
				return filepath.SkipDir
			}
			size, err := calculatePathSize(path)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", path, err))
				return filepath.SkipDir
			}
			ageDays := wholeAgeDays(info.ModTime())
			result.Artifacts = append(result.Artifacts, models.PurgeArtifact{
				Path:     path,
				Type:     artifactType,
				Size:     size,
				AgeDays:  ageDays,
				Selected: ageDays >= purgeMinimumAgeDays,
			})
			result.TotalSize += size
			return filepath.SkipDir
		})
		if walkErr != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", configuredPath, walkErr))
		}
	}

	result.ArtifactCount = len(result.Artifacts)
	emitter.Complete(models.PurgeResult{SpaceFreed: result.TotalSize, RemovedCount: result.ArtifactCount, Errors: result.Errors})
	return result, nil
}

func (s *PurgeService) ExecutePurge(paths []string) (models.PurgeResult, error) {
	if err := VerifyBundledScript(filepath.Join(s.scriptsPath, "bin", "purge.sh")); err != nil {
		return models.PurgeResult{}, err
	}

	emitter := NewRuntimeOperationEmitter(s.ctx, "purge")
	result := models.PurgeResult{}

	for _, rawPath := range paths {
		path := filepath.Clean(expandUserPath(rawPath))
		info, err := os.Stat(path)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", rawPath, err))
			continue
		}
		if !info.IsDir() {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: not a directory", rawPath))
			continue
		}
		if _, ok := purgeArtifactNames[filepath.Base(path)]; !ok {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: not a purge artifact", rawPath))
			continue
		}
		size, _ := calculatePathSize(path)
		if err := os.RemoveAll(path); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", rawPath, err))
			continue
		}
		result.SpaceFreed += size
		result.RemovedCount++
	}

	emitter.Complete(result)
	return result, nil
}

func (s *PurgeService) GetPaths() ([]string, error) {
	path := purgeConfigPath()
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return append([]string{}, purgeDefaultPaths...), nil
		}
		return nil, fmt.Errorf("failed to read purge paths: %w", err)
	}
	defer file.Close()

	var paths []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			paths = append(paths, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read purge paths: %w", err)
	}
	if len(paths) == 0 {
		return append([]string{}, purgeDefaultPaths...), nil
	}
	return paths, nil
}

func (s *PurgeService) UpdatePaths(paths []string) error {
	configPath := purgeConfigPath()
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create purge config directory: %w", err)
	}

	cleaned := make([]string, 0, len(paths))
	for _, path := range paths {
		path = strings.TrimSpace(path)
		if path != "" {
			cleaned = append(cleaned, path)
		}
	}

	tempPath := configPath + ".tmp"
	if err := os.WriteFile(tempPath, []byte(strings.Join(cleaned, "\n")+"\n"), 0644); err != nil {
		return fmt.Errorf("failed to write purge paths: %w", err)
	}
	if err := os.Rename(tempPath, configPath); err != nil {
		_ = os.Remove(tempPath)
		return fmt.Errorf("failed to save purge paths: %w", err)
	}
	return nil
}

func purgeConfigPath() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "mole", "purge_paths")
}

func expandUserPath(path string) string {
	if path == "~" {
		return os.Getenv("HOME")
	}
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(os.Getenv("HOME"), strings.TrimPrefix(path, "~/"))
	}
	return path
}

func wholeAgeDays(modTime time.Time) int {
	days := int(time.Since(modTime).Hours() / 24)
	if days < 0 {
		return 0
	}
	return days
}

func calculatePathSize(path string) (int64, error) {
	var size int64
	err := filepath.WalkDir(path, func(_ string, entry fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return nil
		}
		size += info.Size()
		return nil
	})
	return size, err
}
