#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"/..

mkdir -p bin
cd bin
curl \
  --location \
  --remote-name \
  https://github.com/mikefarah/yq/releases/download/v4.16.1/yq_darwin_amd64
mv yq_darwin_amd64 yq
chmod +x yq
