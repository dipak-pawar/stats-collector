#!/usr/bin/env bash

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker push dipakpawar231/stats-collector:0.1
docker push dipakpawar231/stats-collector:latest
