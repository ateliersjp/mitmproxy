#!/bin/sh

DOCKER_BUILDKIT=1 docker build --build-arg IGNORECACHE=$(date +%s) -t ateliersjp/mitmproxy .
docker push ateliersjp/mitmproxy
