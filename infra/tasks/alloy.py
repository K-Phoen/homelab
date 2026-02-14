from pyinfra.context import host
from pyinfra.api.deploy import deploy
from pyinfra.operations import apt, files, systemd

from .defaults import DEFAULTS
from .utils import resource_path

@deploy("Install Alloy", data_defaults=DEFAULTS)
def install():
    """
    Based on https://grafana.com/docs/alloy/latest/set-up/install/linux/
    """
    files.directory(
        name="Ensure APT keyrings directory exists",
        path="/etc/apt/keyrings",
        user="root",
        group="root",
        present=True,
    )

    files.download(
        name="Download GPG key",
        src="https://apt.grafana.com/gpg-full.key",
        dest="/etc/apt/keyrings/grafana.asc",
        user="root",
        group="root",
        mode="644",
    )

    files.put(
        name="Ensure Grafana source file exists",
        src=resource_path("files/grafana.sources"),
        dest="/etc/apt/sources.list.d/grafana.sources",
        mode="644",
        user="root",
        group="root",
    )

    apt.packages(
        name="Install Alloy package",
        packages=["alloy"],
        no_recommends=True,
        update=True,
        cache_time=3600, # seconds
    )

    alloy_config = files.template(
        name="Configure Alloy",
        src=resource_path("templates/config.alloy.j2"),
        dest="/etc/alloy/config.alloy",
        mode="600",
        user="alloy",
        group="alloy",
        grafana_cloud_prometheus_username=host.data.gcloud_prometheus_username,
        grafana_cloud_loki_username=host.data.gcloud_loki_username,
        grafana_cloud_token=host.data.gcloud_token,
        node_exporter_enabled=host.data.node_exporter_enabled,
        blocky_enabled=host.data.blocky_enabled,
        k3s_master=host.data.k3s_master,
    )

    systemd.service(
        name="Restart and enable the Alloy service",
        service="alloy.service",
        running=True,
        restarted=alloy_config.changed, # Trigger a restart only if the config changed
        enabled=True,
    )
