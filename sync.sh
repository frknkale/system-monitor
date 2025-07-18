#!/bin/bash

SRC="/mnt/d/Downloads/jotform/lfs-shared/monitoring-merged/"

DEST="furkan@34.42.118.72:/opt/monitoring-merged/"

rsync -avz --checksum --delete  --rsync-path="sudo rsync" --exclude='.git/' --exclude='go.mod' --exclude='go.sum' "$SRC" "$DEST"
