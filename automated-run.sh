#!/bin/bash

loop until [ -f "Prompt.md" ]; do
  echo "Waiting for Prompt.md to be created..."
  opencode run --agent yolo --model "omlx1/Qwen3.5-27B-Claude-4.6-Opus-Distilled-MLX-6bit" "$(cat Prompt.md)"
  sleep 300 # Check every 5 minutes
done


