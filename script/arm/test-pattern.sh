#!/usr/bin/env bash

set -euo pipefail

if [ $# -ne 1 ]; then
    echo "Usage: $0 <test_pattern>"
    exit 1
fi

pattern="$1"

docker compose exec strv-newsletter go test -race -v -run "$pattern" ./... | \
sed -E 's/===\s+RUN/=== \x1B[33mRUN\x1B[0m/g' | \
sed -E 's/===/\x1B[36m&\x1B[0m/g' | \
sed -E 's/---/\x1B[35m&\x1B[0m/g' | \
sed -E $'/PASS:/s/(PASS:\\s[^ ]*\\s(\\S*))/\x1B[32m&\x1B[0m/' | \
sed -E $'/FAIL:/s/(FAIL:\\s[^ ]*\\s(\\S*))/\x1B[31m&\x1B[0m/'
