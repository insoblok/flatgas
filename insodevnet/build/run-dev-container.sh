#!/bin/sh
# Run the dev container interactively with shell access

IMAGE_NAME="insoblok/inso-node:dev"

echo "🚀 Launching interactive shell in $IMAGE_NAME"
docker run -it --rm --entrypoint /bin/bash "$IMAGE_NAME"
