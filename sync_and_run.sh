#!/bin/bash

SRC="/mnt/d/Downloads/jotform/lfs-shared/monitoring/"

DEST="furkan@34.42.118.72:/opt/monitoring/"

echo "🔄 Syncing files to remote server..."
rsync -avz --checksum --delete \
  --rsync-path="sudo rsync" \
  --exclude='.git/' \
  --exclude='go.mod' \
  --exclude='go.sum' \
  --exclude='output/' \
  "$SRC" "$DEST"

echo "Building project on remote server..."
ssh furkan@34.42.118.72 << 'EOF'
  echo "Changing to project directory..."
  cd /opt/monitoring || exit 1

  echo "Tidying modules..."
  go mod tidy

  echo "Running..."
  sudo go build -o monitoring-app

  echo "Build completed. Binary: /opt/monitoring/monitoring-app"
EOF
