#!/bin/bash

set -euo pipefail

GREEN="\033[0;32m"
RED="\033[0;31m"
YELLOW="\033[1;33m"
RESET="\033[0m"

SEISCOMP_ROOT="${SEISCOMP_ROOT:-}"

check_seiscomp_root() {
    if [[ -z "$SEISCOMP_ROOT" ]]; then
        echo -e "${RED}Error: SEISCOMP_ROOT environment variable is not set.${RESET}"
        exit 1
    fi
    if [[ ! -d "$SEISCOMP_ROOT" ]]; then
        echo -e "${RED}Error: SEISCOMP_ROOT directory does not exist: $SEISCOMP_ROOT${RESET}"
        exit 1
    fi
}

check_rsync() {
    if ! command -v rsync &>/dev/null; then
        echo -e "${RED}Error: rsync is not installed. Please install rsync and try again.${RESET}"
        exit 1
    fi
}

check_seiscomp_directory() {
    if [[ ! -d "seiscomp" ]]; then
        echo -e "${RED}Error: seiscomp directory does not exist.${RESET}"
        exit 1
    fi
}

install_plugin() {
    echo -e "${YELLOW}Installing AnyShake plugin...${RESET}"
    check_seiscomp_directory
    rsync -av --progress seiscomp/ "$SEISCOMP_ROOT/" || {
        echo -e "${RED}Error: Failed to copy files to $SEISCOMP_ROOT.${RESET}"
        exit 1
    }
    echo -e "${GREEN}AnyShake plugin installed successfully.${RESET}"
}

main() {
    check_seiscomp_root
    check_rsync
    install_plugin
}

main
