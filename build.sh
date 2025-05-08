#!/bin/bash

set -euo pipefail

GREEN="\033[0;32m"
RED="\033[0;31m"
YELLOW="\033[1;33m"
RESET="\033[0m"

PLUGIN_BIN="anyshake_plugin"
PLUGIN_DIR="./anyshake_plugin"
SCRIPTS_DIR="./scripts"
BUILD_ROOT="./build"

usage() {
    echo -e "${YELLOW}Usage: $0 <GOOS> <GOARCH> [GOARM] [GOMIPS] <OUTPUT_PREFIX>${RESET}"
    echo -e "Example: $0 linux amd64 '' '' linux-amd64"
}

check_go_toolchain() {
    if ! command -v go &>/dev/null; then
        echo -e "${RED}Error: Go toolchain not found. Please install Go.${RESET}"
        exit 1
    fi
}

prepare_directories() {
    local output_prefix="$1"
    local build_dir="$BUILD_ROOT/$output_prefix"

    echo -e "${YELLOW}Preparing directories...${RESET}"
    mkdir -p "$build_dir/seiscomp/share/plugins/seedlink/"
    mkdir -p "$build_dir/seiscomp/etc/descriptions/"
    mkdir -p "$build_dir/seiscomp/share/templates/seedlink/"

    # Ensure cleanup on failure
    trap "cleanup_build_dir '$output_prefix'" EXIT
}

build_plugin() {
    local output_prefix="$1"
    local build_dir="$BUILD_ROOT/$output_prefix"
    local output_binary="$build_dir/seiscomp/share/plugins/seedlink/$PLUGIN_BIN"

    echo -e "${YELLOW}Building $PLUGIN_BIN for $GOOS/$GOARCH...${RESET}"
    (cd "$PLUGIN_DIR" && CGO_ENABLED=0 go build -ldflags="-s -w" -v -trimpath -o $PLUGIN_BIN) || {
        echo -e "${RED}Error: Failed to build $PLUGIN_BIN for $GOOS/$GOARCH.${RESET}"
        exit 1
    }

    echo -e "${YELLOW}Copying $PLUGIN_BIN...${RESET}"
    mv -v "$PLUGIN_DIR/$PLUGIN_BIN" "$output_binary" || {
        echo -e "${RED}Error: Failed to copy $PLUGIN_BIN.${RESET}"
        exit 1
    }

    echo -e "${GREEN}Build completed: $output_binary${RESET}"
}

copy_config_assets() {
    local output_prefix="$1"
    local build_dir="$BUILD_ROOT/$output_prefix"

    echo -e "${YELLOW}Copying description files...${RESET}"
    cp -vr "$PLUGIN_DIR/descriptions/"*.xml "$build_dir/seiscomp/etc/descriptions/" || {
        echo -e "${RED}Error: Failed to copy description files.${RESET}"
        exit 1
    }

    echo -e "${YELLOW}Copying template files...${RESET}"
    cp -vr templates/* "$build_dir/seiscomp/share/templates/seedlink/" || {
        echo -e "${RED}Error: Failed to copy template files.${RESET}"
        exit 1
    }
}

copy_setup_scripts() {
    local output_prefix="$1"
    local build_dir="$BUILD_ROOT/$output_prefix"

    echo -e "${YELLOW}Copying setup scripts...${RESET}"
    cp -vr "$SCRIPTS_DIR/"* "$build_dir/" || {
        echo -e "${RED}Error: Failed to copy setup scripts.${RESET}"
        exit 1
    }
}

create_archive() {
    local output_prefix="$1"
    local build_dir="$BUILD_ROOT/$output_prefix"
    local archive_path="$BUILD_ROOT/$output_prefix.tar.gz"

    echo -e "${YELLOW}Creating final archive...${RESET}"
    tar -cvzf "$archive_path" -C "$build_dir" . || {
        echo -e "${RED}Error: Failed to create final archive.${RESET}"
        exit 1
    }

    echo -e "${GREEN}Packaging completed successfully: $archive_path${RESET}"
}

cleanup_build_dir() {
    local output_prefix="$1"
    local build_dir="$BUILD_ROOT/$output_prefix"

    echo -e "${YELLOW}Cleaning up build directory...${RESET}"
    rm -rf "$build_dir" || {
        echo -e "${RED}Error: Failed to cleanup build directory.${RESET}"
        exit 1
    }
    trap - EXIT
}

main() {
    if [[ $# -lt 5 ]]; then
        usage
        exit 1
    fi

    export GOOS="$1"
    export GOARCH="$2"
    export GOARM="${3:-}"
    export GOMIPS="${4:-}"
    local output_prefix="$5"

    check_go_toolchain
    prepare_directories "$output_prefix"
    build_plugin "$output_prefix"
    copy_config_assets "$output_prefix"
    copy_setup_scripts "$output_prefix"
    create_archive "$output_prefix"
    cleanup_build_dir "$output_prefix"

    echo -e "${GREEN}Build completed successfully: $BUILD_ROOT/$output_prefix.tar.gz${RESET}"
}

main "$@"
