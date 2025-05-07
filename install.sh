#!/bin/bash

set -e

GREEN="\033[0;32m"
RED="\033[0;31m"
YELLOW="\033[1;33m"
RESET="\033[0m"

if [[ -z "$SEISCOMP_ROOT" ]]; then
    echo -e "${RED}Error: SEISCOMP_ROOT environment variable is not set.${RESET}"
    exit 1
fi

PLUGIN_DIR="anyshake_plugin"
PLUGIN_DST="$SEISCOMP_ROOT/share/plugins/seedlink/"
DESCRIPTION_SRC="anyshake_plugin/descriptions/*.xml"
DESCRIPTION_DST="$SEISCOMP_ROOT/etc/descriptions/"
TEMPLATE_SRC="templates/*"
TEMPLATE_DST="$SEISCOMP_ROOT/share/templates/seedlink/"

build_plugin() {
    echo -e "${YELLOW}Checking Go toolchain...${RESET}"
    if ! command -v go &>/dev/null; then
        echo -e "${RED}Error: Go toolchain not found. Please install Go before proceeding.${RESET}"
        exit 1
    fi

    echo -e "${YELLOW}Building anyshake_plugin...${RESET}"
    mkdir -p "$PLUGIN_DST"
    (cd "$PLUGIN_DIR" && go build -o "$PLUGIN_DST/anyshake_plugin") || {
        echo -e "${RED}Failed to build anyshake_plugin.${RESET}"
        exit 1
    }
    echo -e "${GREEN}Build successful.${RESET}"
}

install_plugin() {
    echo -e "${YELLOW}Installing AnyShake plugin...${RESET}"

    if [[ ! -f "$PLUGIN_DST/anyshake_plugin" ]]; then
        build_plugin
    fi

    echo -e "${YELLOW}Cleaning up old files...${RESET}"
    rm -f "$DESCRIPTION_DST/anyshake_plugin.xml" 2>/dev/null || true
    rm -rf "$TEMPLATE_DST/anyshake" 2>/dev/null || true

    mkdir -p "$DESCRIPTION_DST" "$TEMPLATE_DST"

    echo -e "${YELLOW}Copying description files...${RESET}"
    cp -v $DESCRIPTION_SRC "$DESCRIPTION_DST" || {
        echo -e "${RED}Failed to copy description files.${RESET}"
        exit 1
    }

    echo -e "${YELLOW}Copying template files...${RESET}"
    cp -vr $TEMPLATE_SRC "$TEMPLATE_DST" || {
        echo -e "${RED}Failed to copy template files.${RESET}"
        exit 1
    }

    echo -e "${GREEN}AnyShake plugin installed successfully.${RESET}"
}

uninstall_plugin() {
    echo -e "${YELLOW}Uninstalling AnyShake plugin...${RESET}"
    rm -f "$PLUGIN_DST/anyshake_plugin" 2>/dev/null || true
    rm -f "$DESCRIPTION_DST/anyshake_plugin.xml" 2>/dev/null || true
    rm -rf "$TEMPLATE_DST/anyshake" 2>/dev/null || true

    echo -e "${GREEN}AnyShake plugin uninstalled successfully.${RESET}"
}

case "$1" in
    install)
        install_plugin
        ;;
    uninstall)
        uninstall_plugin
        ;;
    *)
        echo -e "${RED}Usage: $0 {install|uninstall}${RESET}"
        exit 1
        ;;
esac
