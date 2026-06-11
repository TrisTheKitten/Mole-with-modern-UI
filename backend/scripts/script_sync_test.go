package scripts_test

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

const cliRepoEnvVar = "MOLE_CLI_REPO"

var syncedScriptRelPaths = []string{
	"lib/core/app_protection_data.sh",
	"lib/core/bundle_resolver.sh",
	"lib/core/commands.sh",
	"lib/core/help.sh",
	"lib/core/history.sh",
	"lib/core/pkg_receipts.sh",
	"lib/core/timeouts.sh",
	"bin/installer.sh",
	"bin/history.sh",
	"lib/clean/hints.sh",
	"lib/clean/launch_services.sh",
	"lib/clean/maven.sh",
	"lib/clean/purge_shared.sh",
	"lib/optimize/diagnostics.sh",
	"lib/uninstall/brew.sh",
	"lib/manage/purge_paths.sh",
}

var removedScriptRelPaths = []string{
	"lib/manage/autofix.sh",
	"lib/manage/update.sh",
}

func guiRepoRootFromCaller(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to resolve test source file location")
	}
	return filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))
}

func cliRepoRoot(t *testing.T) string {
	t.Helper()
	repo := os.Getenv(cliRepoEnvVar)
	if repo == "" {
		t.Skipf("%s is not set; skipping bundled-script byte-identity test", cliRepoEnvVar)
	}
	info, err := os.Stat(repo)
	if err != nil || !info.IsDir() {
		t.Skipf("%s=%q does not point to an existing directory; skipping", cliRepoEnvVar, repo)
	}
	return repo
}

func assertByteIdentical(t *testing.T, rel, guiPath, cliPath string) {
	t.Helper()
	guiBytes, err := os.ReadFile(guiPath)
	if err != nil {
		t.Errorf("%s: cannot read bundled script: %v", rel, err)
		return
	}
	cliBytes, err := os.ReadFile(cliPath)
	if err != nil {
		t.Errorf("%s: cannot read CLI counterpart: %v", rel, err)
		return
	}
	if !bytes.Equal(guiBytes, cliBytes) {
		t.Errorf("%s: bundled script differs from CLI counterpart (%d vs %d bytes)", rel, len(guiBytes), len(cliBytes))
	}
}

func TestBundledScriptsByteIdentical(t *testing.T) {
	cliRepo := cliRepoRoot(t)
	guiScriptsDir := filepath.Join(guiRepoRootFromCaller(t), "scripts")

	t.Run("SyncedInventoryPresentAndIdentical", func(t *testing.T) {
		for _, rel := range syncedScriptRelPaths {
			guiPath := filepath.Join(guiScriptsDir, rel)
			cliPath := filepath.Join(cliRepo, rel)
			if _, err := os.Stat(guiPath); err != nil {
				t.Errorf("%s: required bundled script missing: %v", rel, err)
				continue
			}
			if _, err := os.Stat(cliPath); err != nil {
				t.Errorf("%s: CLI counterpart missing: %v", rel, err)
				continue
			}
			assertByteIdentical(t, rel, guiPath, cliPath)
		}
	})

	t.Run("EverySharedRelativePathIdentical", func(t *testing.T) {
		comparedCount := 0
		walkErr := filepath.Walk(guiScriptsDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			rel, relErr := filepath.Rel(guiScriptsDir, path)
			if relErr != nil {
				return relErr
			}
			cliPath := filepath.Join(cliRepo, rel)
			if _, statErr := os.Stat(cliPath); statErr != nil {
				return nil
			}
			comparedCount++
			assertByteIdentical(t, rel, path, cliPath)
			return nil
		})
		if walkErr != nil {
			t.Fatalf("failed to walk bundled scripts directory: %v", walkErr)
		}
		if comparedCount == 0 {
			t.Fatal("found no shared relative-path scripts to compare")
		}
	})

	t.Run("RemovedScriptsAbsent", func(t *testing.T) {
		for _, rel := range removedScriptRelPaths {
			path := filepath.Join(guiScriptsDir, rel)
			if _, err := os.Stat(path); !os.IsNotExist(err) {
				t.Errorf("%s: stale script must be absent but still exists (stat err: %v)", rel, err)
			}
		}
	})
}
