from pyinfra.api.deploy import deploy
from pyinfra.operations import files, server

from .defaults import DEFAULTS
from .utils import resource_path

@deploy("Setup k3s", data_defaults=DEFAULTS)
def install():
    etcd_backups()

@deploy("Setup etcd backups", data_defaults=DEFAULTS)
def etcd_backups():
    backups_dir = "/mnt/k3s-backup"

    files.directory(
        name="Ensure k3s backup mountpoint exists",
        path=backups_dir,
        present=True,
        user="root",
        group="root",
        mode="700",
        recursive=True,
    )

    files.line(
        name="Add NFS volume to fstab",
        path="/etc/fstab",
        line=f"beet.lab:/mnt/main/Backups/k3s {backups_dir} nfs defaults 0 0",
        present=True,
    )

    server.shell(
        name="Mount NFS volume",
        commands=[f"mount {backups_dir}"]
    )

    files.put(
        name="Create k3s config file",
        src=resource_path("files/k3s/k3s-config.yaml"),
        dest="/etc/rancher/k3s/config.yaml",
        user="root",
        group="root",
        mode="644",
    )
