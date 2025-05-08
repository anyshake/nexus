#!/bin/bash

set -euo pipefail

GREEN="\033[0;32m"
RED="\033[0;31m"
YELLOW="\033[1;33m"
RESET="\033[0m"

SEISCOMP_ROOT="${SEISCOMP_ROOT:-}"

check_environment() {
    if [[ -z "$SEISCOMP_ROOT" ]]; then
        echo -e "${RED}Error: SEISCOMP_ROOT environment variable is not set.${RESET}"
        exit 1
    fi
}

remove_file() {
    local file_path="$1"
    if [[ -f "$file_path" ]]; then
        echo -e "${YELLOW}Removing $file_path...${RESET}"
        rm -v "$file_path" || {
            echo -e "${RED}Error: Failed to remove $file_path.${RESET}"
            exit 1
        }
    fi
}

remove_directory() {
    local dir_path="$1"
    if [[ -d "$dir_path" ]]; then
        echo -e "${YELLOW}Removing directory $dir_path...${RESET}"
        rm -rv "$dir_path" || {
            echo -e "${RED}Error: Failed to remove directory $dir_path.${RESET}"
            exit 1
        }
    fi
}

uninstall_plugin() {
    local plugin_path="$SEISCOMP_ROOT/share/plugins/seedlink/anyshake_plugin"
    local description_path="$SEISCOMP_ROOT/etc/descriptions/anyshake_plugin.xml"
    local template_dir="$SEISCOMP_ROOT/share/templates/seedlink/anyshake"

    remove_file "$plugin_path"
    remove_file "$description_path"
    remove_directory "$template_dir"

    echo -e "${GREEN}AnyShake plugin uninstalled successfully.${RESET}"
}

main() {
    check_environment
    uninstall_plugin
}

main "$@"
