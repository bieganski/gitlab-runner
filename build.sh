#!/bin/bash

set -u
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"


pushd $SCRIPT_DIR > /dev/null

sudo make runner-and-helper-bin-host || true

echo "---------------------------------------------"
ls -l out/binaries/gitlab-runner

echo "===== REALPATH: "
realpath out/binaries/gitlab-runner
sudo cp out/binaries/gitlab-runner  /usr/local/ci
realpath /usr/local/ci/gitlab-runner

popd > /dev/null
