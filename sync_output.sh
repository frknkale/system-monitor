#!/bin/bash

REMOTE="furkan@34.42.118.72:/opt/monitoring-merged/"
LOCAL="/mnt/d/Downloads/jotform/lfs-shared/monitoring-merged/"

rsync -avz --delete --rsync-path="sudo rsync" --exclude='.git/' --exclude='go.mod' --exclude='go.sum' "$REMOTE/output/" "$LOCAL/output/"

rsync -avz --delete --rsync-path="sudo rsync" --exclude='.git/' --exclude='go.mod' --exclude='go.sum' "$REMOTE/logs/" "$LOCAL/logs/"