#!/bin/bash

set -u
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"


pushd $SCRIPT_DIR > /dev/null

sudo make runner-and-helper-bin-host

echo "---------------------------------------------"
ls -l out/binaries/gitlab-runner

echo "===== REALPATH: "
realpath out/binaries/gitlab-runner

popd > /dev/null
