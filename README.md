# MoleUI

A native macOS app for cleaning, optimizing, and monitoring your Mac. MoleUI wraps the [Mole CLI](https://github.com/tw93/Mole) (v1.42.0 Otter) in a graphical interface so you can run system maintenance without the terminal.

## Disclaimer !

**MoleUI is not affiliated with, endorsed by, or maintained by [tw93](https://github.com/tw93)** — the author of [Mole](https://github.com/tw93/Mole). This is an independent, community-driven GUI built on top of the open-source Mole CLI. It is not the official Mole app.

If you want the native Mac experience that tw93 ships and maintains, use **[Mole for Mac](https://mole.fit)** — the official paid GUI. Buying it directly supports his work on Mole and helps keep the free CLI moving forward. We recommend it over this project when you can.

Because MoleUI is maintained separately, it can fall behind the official CLI and app on features, fixes, and polish. Some capabilities may be incomplete or arrive later (for example, History, Purge, and Touch ID are still pending here). Use this repo if you specifically want a free, open-source desktop UI — not as a drop-in replacement for the official product.

## What you can do

- **Clean** — Scan and remove caches, logs, and temporary files to free disk space
- **Optimize** — Run maintenance tasks that rebuild caches and refresh system services
- **Installer Cleanup** — Find and remove leftover installer downloads (`.dmg`, `.pkg`, `.zip`, and more)
- **Uninstall** — Remove apps and their leftover files, with dry-run preview
- **History** — Review past Mole sessions, deletions, and operation logs *(coming soon)*
- **Analyze** — Browse disk usage by folder, scan external volumes, and spot large files
- **Status** — Watch CPU, memory, disk, GPU, network, battery, and system health in real time
- **Purge** — Remove stale build artifacts from dev folders *(coming soon)*
- **Touch ID** — Set up Touch ID for `sudo` *(coming soon)*

## Credits

**MoleUI** was created by [Roniel](https://github.com/roniel-rhack) ([@roniel-rhack](https://github.com/roniel-rhack)). The original desktop app, Wails integration, and Vue frontend live at [github.com/roniel-rhack/mole-ui](https://github.com/roniel-rhack/mole-ui).

This fork is maintained by [Tris The Kitten](https://github.com/TrisTheKitten) with UI updates and Mole CLI feature parity work.

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
   git clone https://github.com/TrisTheKitten/Mole-with-modern-UI.git
   cd Mole-with-modern-UI
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
make build
```

Output: `build/bin/mole-wails.app`

`make build` runs a clean `wails build`, then ad-hoc signs the bundle and strips the macOS quarantine attribute so the app opens on double-click without a separate terminal step. Plain `wails build` still works, but a copy that gets zipped, moved, or downloaded may trigger a Gatekeeper "damaged or incomplete" warning because the app is ad-hoc signed rather than notarized (notarization requires a paid Apple Developer account). If that happens, run:

```bash
xattr -cr build/bin/mole-wails.app
```

## Tech stack

| Layer | Tools |
| --- | --- |
| Backend | Go 1.24+, Wails v2, macOS APIs |
| Frontend | Vue 3, Pinia, TypeScript, Vite, PrimeIcons |
| Engine | Bash scripts and Go modules from Mole CLI v1.42.0 |

**Architecture:** Wails events push live progress to the UI. Go services handle system work; Vue components handle display and input.

## Project structure

```
mole-wails/
├── backend/                 # Go backend
│   ├── services/           # Clean, uninstall, optimize, installer, history, purge, Touch ID
│   ├── models/             # Shared data types
│   ├── analyze/            # Disk analyzer (from Mole)
│   └── status/             # System monitor and process watch (from Mole)
├── frontend/               # Vue 3 app
│   └── src/
│       ├── components/
│       │   ├── tabs/       # Feature screens (Clean, Optimize, Installer, History, etc.)
│       │   ├── layout/     # Sidebar and shell
│       │   ├── status/     # Health summary and resource sections
│       │   └── shared/     # Buttons, dialogs, panels, and shared UI
│       └── stores/         # Pinia state per feature tab
├── scripts/                # Bash scripts from Mole CLI
│   ├── bin/                # Entry-point scripts (clean, uninstall, optimize, history, etc.)
│   └── lib/                # Shared shell libraries
└── AGENT.md                # Architecture and safety guidelines for contributors
```

## What's included

| Area | Status |
| --- | --- |
| Clean | Scan, review, and remove caches with progress tracking and whitelist support |
| Optimize | Maintenance tasks with dry-run preview and task whitelist |
| Installer Cleanup | Scan and remove leftover installer files from common download locations |
| Uninstall | Remove apps and associated files with dry-run preview |
| History | Backend ready; UI coming soon |
| Analyze | Folder breakdown, large-file detection, external volume scanning, Finder integration |
| Status | Live CPU, memory, disk, GPU, network, battery, health score, and process alerts |
| Purge | Backend ready; UI coming soon |
| Touch ID | Backend ready; UI coming soon |
| About | Version info, feature overview, and credits |
| Dialogs | Custom confirm flows (browser `alert`/`confirm` do not work in Wails) |
| Progress | Real-time updates over Wails events |
| Errors | Actionable messages in the UI |

### UI details

- Dark theme
- Loading states and progress bars
- Toast notifications
- Confirmation before destructive actions
- Dry-run previews for uninstall and optimize

## For contributors

### Dialogs and notifications

Wails desktop apps cannot use `alert()`, `confirm()`, or `prompt()`. Use project components instead:

- `ConfirmDialog.vue` for confirmations
- `Toast.vue` for short feedback
- `MessageBanner.vue` for inline errors
- `EmptyState.vue` for empty views

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

Feature tabs with dedicated stores:

- `useCleanStore()`
- `useOptimizeStore()`
- `useInstallerStore()`
- `useUninstallStore()`
- `useHistoryStore()`
- `usePurgeStore()` *(backend wired; UI not yet exposed)*
- `useTouchIDStore()` *(backend wired; UI not yet exposed)*

Analyze and Status use local component state and Wails bindings directly.

## Known limitations

- **Cloud and sparse files:** Size may show as 0 for iCloud or sparse files. The app falls back to logical size when needed.
- **Disk metrics:** Reports the primary disk only to avoid double-counting APFS containers.
- **Purge and Touch ID:** Backend services exist but the sidebar entries are disabled until the UI ships.

## Contributing

Pull requests are welcome. Read `AGENT.md` for architecture and safety guidelines before changing cleanup or uninstall behavior.

## License

MIT
