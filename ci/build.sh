#!/usr/bin/env bash

docker build -t "timothylock/inventory-management:latest" .
docker push timothylock/inventory-management:latest