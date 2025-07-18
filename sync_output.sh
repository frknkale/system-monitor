#!/bin/bash

REMOTE="furkan@34.42.118.72:/opt/monitoring-merged/output/"
LOCAL="/mnt/d/Downloads/jotform/lfs-shared/monitoring-merged/output/"

rsync -avz --delete --rsync-path="sudo rsync" "$REMOTE" "$LOCAL"