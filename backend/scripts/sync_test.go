package scripts_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

const cliSourceEnvVar = "MOLE_CLI_SOURCE"

const defaultCLISiblingDir = "Mole-main"

var sharedScriptPaths = []string{
	"bin/installer.sh",
	"bin/history.sh",
	"lib/core/app_protection_data.sh",
	"lib/core/bundle_resolver.sh",
	"lib/core/commands.sh",
	"lib/core/help.sh",
	"lib/core/history.sh",
	"lib/core/pkg_receipts.sh",
	"lib/core/timeouts.sh",
	"lib/clean/hints.sh",
	"lib/clean/launch_services.sh",
	"lib/clean/maven.sh",
	"lib/clean/purge_shared.sh",
	"lib/optimize/diagnostics.sh",
	"lib/uninstall/brew.sh",
	"lib/manage/purge_paths.sh",
}

var removedScriptPaths = []string{
	"lib/manage/autofix.sh",
	"lib/manage/update.sh",
}

func guiRepoRoot(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("resolve working directory: %v", err)
	}

	for {
		if _, statErr := os.Stat(filepath.Join(dir, "go.mod")); statErr == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("could not locate GUI repo root containing go.mod")
		}
		dir = parent
	}
}

func cliSourceRoot(t *testing.T, guiRoot string) string {
	t.Helper()

	if override := os.Getenv(cliSourceEnvVar); override != "" {
		return override
	}
	return filepath.Join(filepath.Dir(guiRoot), defaultCLISiblingDir)
}

func TestBundledScriptsAreByteIdenticalToCLI(t *testing.T) {
	guiRoot := guiRepoRoot(t)
	cliRoot := cliSourceRoot(t, guiRoot)

	if info, err := os.Stat(cliRoot); err != nil || !info.IsDir() {
		t.Skipf("CLI source tree not present at %s; set %s to run this test", cliRoot, cliSourceEnvVar)
	}

	guiScriptsRoot := filepath.Join(guiRoot, "scripts")

	for _, relativePath := range sharedScriptPaths {
		bundledPath := filepath.Join(guiScriptsRoot, relativePath)
		cliPath := filepath.Join(cliRoot, relativePath)

		bundledBytes, err := os.ReadFile(bundledPath)
		if err != nil {
			t.Errorf("read bundled script %s: %v", relativePath, err)
			continue
		}

		cliBytes, err := os.ReadFile(cliPath)
		if err != nil {
			t.Errorf("read CLI script %s: %v", relativePath, err)
			continue
		}

		if !bytes.Equal(bundledBytes, cliBytes) {
			t.Errorf("bundled script %s differs from CLI counterpart", relativePath)
		}
	}
}

func TestRemovedManageScriptsAreAbsent(t *testing.T) {
	guiScriptsRoot := filepath.Join(guiRepoRoot(t), "scripts")

	for _, relativePath := range removedScriptPaths {
		fullPath := filepath.Join(guiScriptsRoot, relativePath)
		if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
			t.Errorf("expected removed script %s to be absent, but it is present", relativePath)
		}
	}
}
