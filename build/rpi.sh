#!/bin/bash

echo "Start building the tgbot executable for RPI..."
GOOS=linux GOARCH=arm go build ./cmd/tgbot
echo "Done."