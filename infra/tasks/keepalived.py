from pyinfra.context import host
from pyinfra.api.deploy import deploy
from pyinfra.facts.server import Command
from pyinfra.operations import apt, files, systemd

from .defaults import DEFAULTS
from .utils import resource_path

@deploy("Install keepalived", data_defaults=DEFAULTS)
def install():
    if not (host.data.blocky_enabled or host.data.k3s_master):
        return

    apt.packages(
        name="Install keepalived package",
        packages=["keepalived"],
        no_recommends=True,
        update=True,
        cache_time=3600, # seconds
    )

    default_ipv4_iface = host.get_fact(Command, "ip route get '8.8.8.8' | grep -Po '(?<=(dev ))(\\S+)'")

    config = files.template(
        name="Configure keepalived",
        src=resource_path("templates/keepalived/keepalived.conf.j2"),
        dest="/etc/keepalived/keepalived.conf",
        mode="640",
        user="root",
        group="root",

        blocky_enabled=host.data.blocky_enabled,
        k3s_master=host.data.k3s_master,

        ipv4_iface=default_ipv4_iface,

        # blocky
        vip_blocky=host.data.keepalived_vip_blocky,
        vrrp_priority_blocky=host.data.keepalived_vrrp_priority_blocky,
        virtual_router_id_blocky=host.data.keepalived_virtual_router_id_blocky,

        # k3s
        vip_k3s=host.data.keepalived_vip_k3s,
        vrrp_priority_k3s=host.data.keepalived_vrrp_priority_k3s,
        virtual_router_id_k3s=host.data.keepalived_virtual_router_id_k3s,
    )

    systemd.service(
        name="Restart and enable the keepalived service",
        service="keepalived.service",
        running=True,
        restarted=config.changed, # Trigger a restart only if the config changed
        enabled=True,
    )
