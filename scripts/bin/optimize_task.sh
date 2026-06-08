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
source "$SCRIPT_DIR/lib/optimize/tasks.sh"

ensure_sudo_session "System optimization requires admin access" || true
execute_optimization "$action" ""
