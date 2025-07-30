#!/bin/bash

SRC="/mnt/d/Downloads/jotform/lfs-shared/monitoring/"

DEST="furkan@34.42.118.72:/opt/monitoring/"

rsync -avz --checksum --delete  --rsync-path="sudo rsync" --exclude='.git/' --exclude='go.mod' --exclude='go.sum' --exclude='output/' "$SRC" "$DEST"
