#!/bin/bash
# set -x

readonly WORK_DIR=$(dirname $(readlink --canonicalize-existing "${0}"))

find "${WORK_DIR}" -mindepth 1 -maxdepth 1 -type d -not -name '.git' -exec go test {} \;

exit 0
