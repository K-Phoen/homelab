#!/bin/env bash

# Exit on error. Append "|| true" if you expect an error.
set -o errexit
# Exit on error inside any functions or subshells.
set -o errtrace
# Do not allow use of undefined vars. Use ${VAR:-} to use an undefined VAR
set -o nounset
# Catch the error in case mysqldump fails (but gzip succeeds) in `mysqldump |gzip`
set -o pipefail

cp /var/lib/rancher/k3s/server/token /mnt/k3s-backup/
cp -R /var/lib/rancher/k3s/server/db/ /mnt/k3s-backup/