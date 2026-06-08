# MoleUI

A native macOS app for cleaning, optimizing, and monitoring your Mac. MoleUI wraps the [Mole CLI](https://github.com/tw93/Mole) in a graphical interface so you can run system maintenance without the terminal.

## What you can do

- **Clean** — Scan and remove caches, logs, and temporary files to free disk space
- **Uninstall** — Remove apps and their leftover files in one pass
- **Optimize** — Run maintenance tasks that rebuild caches and refresh system services
- **Analyze** — Browse disk usage by folder and spot large files
- **Status** — Watch CPU, memory, disk, GPU, network, and battery in real time
- **Touch ID** — Set up Touch ID for `sudo` so fewer password prompts interrupt your work

## Credits

**MoleUI** was created by [Roniel](https://github.com/roniel-rhack) ([@roniel-rhack](https://github.com/roniel-rhack)). The original desktop app, Wails integration, and Vue frontend live at [github.com/roniel-rhack/mole-ui](https://github.com/roniel-rhack/mole-ui).

Cleanup, uninstall, and optimization run through [Mole](https://github.com/tw93/Mole) by [@tw93](https://github.com/tw93). MoleUI is the desktop shell; Mole is the engine.

## Get started

### Requirements

- macOS 11.0 or later
- Go 1.24+
- Node.js 16+
- [Wails CLI](https://wails.io/docs/gettingstarted/installation):

  ```bash
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```

### Run locally

1. Clone the repo and open the project folder:

   ```bash
   git clone https://github.com/roniel-rhack/mole-ui
   cd mole-ui
   ```

2. Install dependencies:

   ```bash
   go mod download
   cd frontend && npm install && cd ..
   ```

3. Start the dev build:

   ```bash
   wails dev
   ```

   The app opens with hot reload for frontend and backend changes.

### Build for release

```bash
wails build
```

Output: `build/bin/mole-wails.app`

## Tech stack

| Layer | Tools |
| --- | --- |
| Backend | Go 1.24+, Wails v2, macOS APIs |
| Frontend | Vue 3, Pinia, TypeScript, Vite |
| Engine | Bash scripts and Go modules from Mole |

**Architecture:** Wails events push live progress to the UI. Go services handle system work; Vue components handle display and input.

## Project structure

```
mole-wails/
├── backend/                 # Go backend
│   ├── services/           # Clean, uninstall, optimize, Touch ID
│   ├── models/             # Shared data types
│   ├── analyze/            # Disk analyzer (from Mole)
│   └── status/             # System monitor (from Mole)
├── frontend/               # Vue 3 app
│   └── src/
│       ├── components/
│       │   ├── tabs/       # Feature screens
│       │   ├── layout/     # Sidebar and shell
│       │   └── shared/     # Toast, ConfirmDialog, and shared UI
│       └── stores/         # Pinia state per tab
├── scripts/                # Bash scripts from Mole CLI
└── docs/                   # Plans and reference docs
```

## What's included

| Area | Status |
| --- | --- |
| Clean | Scan, review, and remove caches with progress tracking |
| Uninstall | Remove apps and associated files |
| Optimize | Maintenance tasks with `sudo` support |
| Analyze | Folder breakdown and large-file detection |
| Status | Live CPU, memory, disk, GPU, network, and battery |
| Touch ID | Enable or disable Touch ID for `sudo` |
| About | Version info and credits |
| Dialogs | Custom confirm flows (browser `alert`/`confirm` do not work in Wails) |
| Progress | Real-time updates over Wails events |
| Errors | Actionable messages in the UI |

### UI details

- Dark theme
- Loading states and progress bars
- Toast notifications
- Confirmation before destructive actions

## For contributors

### Dialogs and notifications

Wails desktop apps cannot use `alert()`, `confirm()`, or `prompt()`. Use project components instead:

- `ConfirmDialog.vue` for confirmations
- `Toast.vue` for short feedback
- Inline UI for errors (not `alert()`)

### Wails events

```javascript
// Listen (frontend)
EventsOn('event-name', (data) => {
  console.log('Received:', data)
})
```

```go
// Emit (Go)
runtime.EventsEmit(ctx, "event-name", data)
```

### Pinia stores

Each tab has its own store:

- `useCleanStore()`
- `useUninstallStore()`
- `useOptimizeStore()`
- `useAnalyzeStore()`
- `useStatusStore()`
- `useTouchIDStore()`

## Known limitations

- **Cloud and sparse files:** Size may show as 0 for iCloud or sparse files. The app falls back to logical size when needed.
- **Disk metrics:** Reports the primary disk only to avoid double-counting APFS containers.

## Contributing

Pull requests are welcome on [roniel-rhack/mole-ui](https://github.com/roniel-rhack/mole-ui). Read `AGENT.md` for architecture and safety guidelines before changing cleanup or uninstall behavior.

## License

MIT
