package services

import (
	"fmt"
	"os"
)

// VerifyBundledScript verifies that a bundled script exists and is readable
// before a shell-out service invokes it. It performs read-only checks and makes
// no system changes; on any failure it returns an error so the caller can abort
// without side effects.
func VerifyBundledScript(scriptPath string) error {
	info, err := os.Stat(scriptPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("bundled script not found: %s", scriptPath)
		}
		return fmt.Errorf("bundled script not accessible: %s: %w", scriptPath, err)
	}

	if info.IsDir() {
		return fmt.Errorf("bundled script path is a directory: %s", scriptPath)
	}

	file, err := os.Open(scriptPath)
	if err != nil {
		return fmt.Errorf("bundled script not readable: %s: %w", scriptPath, err)
	}
	defer file.Close()

	return nil
}
