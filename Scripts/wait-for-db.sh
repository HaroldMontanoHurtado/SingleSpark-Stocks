#!/usr/bin/env bash
set -e
host="${1:-cockroach}"
port="${2:-26257}"
echo "Waiting for ${host}:${port}..."
until nc -z "$host" "$port"; do
    sleep 1
done
echo "DB up"
