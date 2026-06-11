#!/bin/bash

set -euo pipefail

export LC_ALL=C
export LANG=C

if [[ $# -lt 1 ]]; then
    echo "Usage: optimize_task.sh <action>" >&2
    exit 1
fi

action="$1"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

source "$SCRIPT_DIR/lib/core/common.sh"
source "$SCRIPT_DIR/lib/core/sudo.sh"
source "$SCRIPT_DIR/lib/optimize/maintenance.sh"
source "$SCRIPT_DIR/lib/optimize/diagnostics.sh"
source "$SCRIPT_DIR/lib/optimize/tasks.sh"

optimize_task_requires_sudo() {
    case "$1" in
        system_maintenance | maintenance_scripts | radio_refresh | swap_cleanup | startup_cache | local_snapshots | network_optimization)
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

if optimize_task_requires_sudo "$action" && ! sudo -n true 2> /dev/null; then
    echo "Admin access required for ${action}. Run a sudo-authenticated action first and retry." >&2
    exit 1
fi

execute_optimization "$action" ""
