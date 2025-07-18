#!/bin/bash

SRC="/mnt/d/Downloads/jotform/lfs-shared/monitoring/"

DEST="crastle@104.154.77.157:/home/crastle/monitoring/"

rsync -avz --checksum --delete  --exclude='.git/' --exclude='go.mod' --exclude='go.sum' "$SRC" "$DEST"
