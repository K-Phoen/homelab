# Backups

## k3s

See https://docs.k3s.io/datastore/backup-restore

The cluster token (`/var/lib/rancher/k3s/server/token`) and SQLite database (`/var/lib/rancher/k3s/server/db/`) are backed up by an ansible-defined cronjob running a backup script every night.

* backup script: `/root/k3s-backup.sh`
* NAS backup dataset: `/mnt/main/Backups/k3s`

## longhorn

TODO

## NAS

TODO
