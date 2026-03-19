from pyinfra.context import host
from pyinfra.api.deploy import deploy
from pyinfra.operations import files, server

from .utils import resource_path

DEFAULTS = {
    "k3s_etcd_backups_nfs": "beet.lab:/mnt/main/Backups/k3s",
    "k3s_etcd_backups_dir": "/mnt/k3s-backup",
}

@deploy("Setup k3s", data_defaults=DEFAULTS)
def install():
    etcd_backups()

@deploy("Setup etcd backups", data_defaults=DEFAULTS)
def etcd_backups():
    files.directory(
        name="Ensure k3s backup mountpoint exists",
        path=host.data.k3s_etcd_backups_dir,
        present=True,
        user="root",
        group="root",
        mode="700",
        recursive=True,
    )

    files.line(
        name="Add NFS volume to fstab",
        path="/etc/fstab",
        line=f"{host.data.k3s_etcd_backups_nfs} {host.data.k3s_etcd_backups_dir} nfs defaults 0 0",
        present=True,
    )

    server.shell(
        name="Mount NFS volume",
        commands=[f"mount {host.data.k3s_etcd_backups_dir}"]
    )

    files.put(
        name="Create k3s config file",
        src=resource_path("files/k3s/k3s-config.yaml"),
        dest="/etc/rancher/k3s/config.yaml",
        user="root",
        group="root",
        mode="644",
    )
