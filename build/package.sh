#!/bin/bash

APP_NAME="GarageBand Launcher"
APP_DIR="build/$APP_NAME.app"
CONTENTS_DIR="$APP_DIR/Contents"
MACOS_DIR="$CONTENTS_DIR/MacOS"
RESOURCES_DIR="$CONTENTS_DIR/Resources"

# Create app structure
mkdir -p "$MACOS_DIR" "$RESOURCES_DIR"

# Build the Go binary
go build -o "$MACOS_DIR/garageband-launcher" ./cmd/garageband-launcher

# Copy Info.plist
cp build/darwin/Info.plist "$CONTENTS_DIR/"

# Create zip for distribution
cd build
zip -r "$APP_NAME.app.zip" "$APP_NAME.app" 