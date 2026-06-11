package main

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestGetScriptsPathDoesNotDependOnWorkingDirectory(t *testing.T) {
	originalWorkingDir := t.TempDir()
	t.Chdir(originalWorkingDir)

	scriptsPath := getScriptsPath()
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve test source path")
	}
	expectedPath := filepath.Join(filepath.Dir(sourceFile), scriptsDirName)

	if scriptsPath != expectedPath {
		t.Fatalf("scripts path = %q, want %q", scriptsPath, expectedPath)
	}
}

func TestHistoryGetDoesNotDependOnWorkingDirectory(t *testing.T) {
	t.Chdir(t.TempDir())

	app := NewApp()
	result, err := app.History.GetHistory(20)
	if err != nil {
		t.Fatalf("history get failed: %v", err)
	}

	if result.Limit != 20 {
		t.Fatalf("history limit = %d, want 20", result.Limit)
	}
}
