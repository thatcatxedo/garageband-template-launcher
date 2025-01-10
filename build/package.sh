#!/bin/bash

set -e  # Exit on any error

APP_NAME="GarageBand Launcher"
APP_DIR="build/$APP_NAME.app"
CONTENTS_DIR="$APP_DIR/Contents"
MACOS_DIR="$CONTENTS_DIR/MacOS"
RESOURCES_DIR="$CONTENTS_DIR/Resources"

echo "Cleaning previous build..."
rm -rf "$APP_DIR"

echo "Creating app structure..."
mkdir -p "$MACOS_DIR" "$RESOURCES_DIR"

echo "Building Go binary..."
# Remove CGO_ENABLED=0 and add required C flags for OpenGL
export CGO_ENABLED=1
export CGO_CFLAGS="-mmacosx-version-min=10.13"
export CGO_LDFLAGS="-mmacosx-version-min=10.13"

go build -o "$MACOS_DIR/garageband-launcher" ./cmd/garageband-launcher

echo "Setting executable permissions..."
chmod +x "$MACOS_DIR/garageband-launcher"

echo "Copying Info.plist..."
cp build/darwin/Info.plist "$CONTENTS_DIR/"

echo "Creating zip for distribution..."
cd build
zip -r "$APP_NAME.app.zip" "$APP_NAME.app"

echo "Build complete!"
echo "App location: $APP_DIR" 