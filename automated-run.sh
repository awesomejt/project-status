#!/bin/bash

while true; do
  if [ -f "status.yaml" ]; then
    status=$(yq -r '.status' status.yaml 2>/dev/null)
    if [[ "$status" == "blocked" || "$status" == "paused" || "$status" == "stopped" ]]; then
      echo "[$(date -Iseconds)] Status is '${status}'. Skipping run..."
      sleep 300
      continue
    fi
  fi

  if [ ! -f "Prompt.md" ]; then
    echo "[$(date -Iseconds)] Waiting for Prompt.md to be created..."
  else
    echo "[$(date -Iseconds)] Running OpenCode..."
    opencode run --agent yolo --model "omlx1/Qwen3.5-27B-Claude-4.6-Opus-Distilled-MLX-6bit" "$(cat Prompt.md)"
  fi
  
  sleep 300 # Run every 5 minutes
done