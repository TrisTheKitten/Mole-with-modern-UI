package services

import (
	"archive/zip"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mole-wails/backend/models"
)

const (
	installerProgressInterval = 500 * time.Millisecond
	installerScanMaxDepth     = 2
	installerZipEntryLimit    = 50
)

var installerExtensions = map[string]bool{
	".dmg":  true,
	".pkg":  true,
	".mpkg": true,
	".iso":  true,
	".xip":  true,
	".zip":  true,
}

type InstallerService struct {
	scriptsPath string
	ctx         context.Context
}

func NewInstallerService(scriptsPath string) *InstallerService {
	return &InstallerService{scriptsPath: scriptsPath}
}

func (s *InstallerService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *InstallerService) ScanInstallers() (models.InstallerScanResult, error) {
	roots := installerScanRoots()
	emitter := NewRuntimeOperationEmitter(s.ctx, "installer")
	result := models.InstallerScanResult{}
	scanned := 0
	lastProgress := time.Time{}

	for _, root := range roots {
		info, err := os.Stat(root.path)
		if err != nil {
			if !os.IsNotExist(err) {
				result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", root.path, err))
			}
			continue
		}
		if !info.IsDir() {
			continue
		}

		walkErr := filepath.WalkDir(root.path, func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", path, walkErr))
				if entry != nil && entry.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			scanned++
			now := time.Now()
			if lastProgress.IsZero() || now.Sub(lastProgress) >= installerProgressInterval {
				emitter.Progress(models.InstallerProgress{CurrentPath: path, ItemsScanned: scanned})
				lastProgress = now
			}
			if entry.IsDir() {
				if exceedsInstallerScanDepth(root.path, path) {
					return filepath.SkipDir
				}
				return nil
			}
			if !isInstallerCandidatePath(path, entry) {
				return nil
			}
			info, err := entry.Info()
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", path, err))
				return nil
			}
			result.Files = append(result.Files, models.InstallerFile{
				Path:         path,
				Size:         info.Size(),
				LastModified: info.ModTime(),
				Source:       root.source,
			})
			result.TotalSize += info.Size()
			return nil
		})
		if walkErr != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", root.path, walkErr))
		}
	}

	if scanned > 0 {
		emitter.Progress(models.InstallerProgress{CurrentPath: "", ItemsScanned: scanned})
	}
	result.FileCount = len(result.Files)
	emitter.Complete(models.InstallerResult{SpaceFreed: result.TotalSize, RemovedCount: result.FileCount, Errors: result.Errors})
	return result, nil
}

func (s *InstallerService) RemoveInstallers(paths []string) (models.InstallerResult, error) {
	if err := VerifyBundledScript(filepath.Join(s.scriptsPath, "bin", "installer.sh")); err != nil {
		return models.InstallerResult{}, err
	}

	emitter := NewRuntimeOperationEmitter(s.ctx, "installer")
	result := models.InstallerResult{}
	for _, rawPath := range paths {
		path := filepath.Clean(expandUserPath(rawPath))
		info, err := os.Lstat(path)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", rawPath, err))
			continue
		}
		if info.Mode()&os.ModeSymlink != 0 {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: symlinks are not supported", rawPath))
			continue
		}
		if info.IsDir() {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: not a file", rawPath))
			continue
		}
		if !isInstallerCandidateFile(path) {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: not a recognized installer", rawPath))
			continue
		}
		if err := os.Remove(path); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", rawPath, err))
			continue
		}
		result.SpaceFreed += info.Size()
		result.RemovedCount++
	}
	emitter.Complete(result)
	return result, nil
}

type installerRoot struct {
	path   string
	source models.InstallerSource
}

func installerScanRoots() []installerRoot {
	home := os.Getenv("HOME")
	return []installerRoot{
		{filepath.Join(home, "Downloads"), models.InstallerSourceDownloads},
		{filepath.Join(home, "Desktop"), models.InstallerSourceDesktop},
		{filepath.Join(home, "Documents"), models.InstallerSourceDocuments},
		{filepath.Join(home, "Public"), models.InstallerSourcePublic},
		{filepath.Join(home, "Library", "Downloads"), models.InstallerSourceLibrary},
		{"/Users/Shared", models.InstallerSourceShared},
		{"/Users/Shared/Downloads", models.InstallerSourceShared},
		{filepath.Join(home, "Library", "Caches", "Homebrew"), models.InstallerSourceHomebrew},
		{filepath.Join(home, "Library", "Mobile Documents", "com~apple~CloudDocs", "Downloads"), models.InstallerSourceICloud},
		{filepath.Join(home, "Library", "Containers", "com.apple.mail", "Data", "Library", "Mail Downloads"), models.InstallerSourceMail},
		{filepath.Join(home, "Library", "Application Support", "Telegram Desktop"), models.InstallerSourceTelegram},
		{filepath.Join(home, "Downloads", "Telegram Desktop"), models.InstallerSourceTelegram},
	}
}

func exceedsInstallerScanDepth(root string, path string) bool {
	if root == path {
		return false
	}
	relativePath, err := filepath.Rel(root, path)
	if err != nil || relativePath == "." {
		return false
	}
	depth := len(strings.Split(relativePath, string(os.PathSeparator)))
	return depth >= installerScanMaxDepth
}

func isInstallerCandidatePath(path string, entry fs.DirEntry) bool {
	if entry.Type()&os.ModeSymlink != 0 {
		return false
	}
	return isInstallerCandidateFile(path)
}

func isInstallerCandidateFile(path string) bool {
	extension := strings.ToLower(filepath.Ext(path))
	if !installerExtensions[extension] {
		return false
	}
	if extension == ".zip" {
		return zipContainsInstaller(path)
	}
	return true
}

func zipContainsInstaller(path string) bool {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return false
	}
	defer reader.Close()

	entryLimit := installerZipEntryLimit
	if len(reader.File) < entryLimit {
		entryLimit = len(reader.File)
	}
	for index := 0; index < entryLimit; index++ {
		entryName := strings.ToLower(reader.File[index].Name)
		for _, extension := range []string{".app", ".pkg", ".dmg", ".xip"} {
			if strings.HasSuffix(entryName, extension) || strings.Contains(entryName, extension+"/") {
				return true
			}
		}
	}
	return false
}
