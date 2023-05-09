#!/bin/bash
docker buildx build --platform linux/amd64,linux/arm64 -t louvandtech/xkcd-bot --push .