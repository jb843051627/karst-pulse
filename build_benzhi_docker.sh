#!/usr/bin/env bash
set -eu

NAME="${1:?image name is required}"
PLATFORM="${2:?platform is required}"
IMAGE="benzhi/${NAME}:latest"
docker build --platform "${PLATFORM}" -f benzhi.Dockerfile -t "${IMAGE}" .
