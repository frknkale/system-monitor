#!/bin/bash

REMOTE="crastle@104.154.77.157:/home/crastle/monitoring/output/"
LOCAL="/mnt/d/Downloads/jotform/lfs-shared/monitoring/output/"

rsync -avz --delete "$REMOTE" "$LOCAL"