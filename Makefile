APP_NAME := mole-wails
APP_BUNDLE := build/bin/$(APP_NAME).app

.PHONY: build dev clean

build: clean
	wails build
	xattr -cr "$(APP_BUNDLE)"
	codesign --force --deep --sign - "$(APP_BUNDLE)"
	@echo "Built and signed $(APP_BUNDLE) — double-click to run, no terminal step needed."

dev:
	wails dev

clean:
	rm -rf build/bin
