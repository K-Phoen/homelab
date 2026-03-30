from pyinfra.context import host
from pyinfra.api.deploy import deploy
from pyinfra.facts.files import File
from pyinfra.facts.server import Arch, Command, Os
from pyinfra.operations import apt, files, server, systemd

from .utils import resource_path

DEFAULTS = {
    "keepalived_vip_blocky": "10.10.10.211",

    # arbitrary unique number from 1 to 255
    # used to differentiate multiple instances of vrrpd
    "keepalived_virtual_router_id_blocky": 42,

    # for electing MASTER, highest priority wins.
    # The valid range of values for priority is [1-255], with priority
    # 255 meaning "address owner".
    # To be MASTER, it is recommended to make this 50 more than on
    # other machines.
    "keepalived_vrrp_priority_blocky": 100,

    "keepalived_vip_k3s": '10.10.10.212',
    "keepalived_vrrp_priority_k3s": 100,
    "keepalived_virtual_router_id_k3s": 43,

    # Exporter
    "keepalived_exporter_version": "1.7.0",
    "keepalived_exporter_dir": "/opt/keepalived-exporter",
    "keepalived_exporter_base_url": "https://github.com/mehdy/keepalived-exporter/releases/download",
    "keepalived_exporter_tmp_dir": "/opt/keepalived-exporter",
}

@deploy("Install keepalived", data_defaults=DEFAULTS)
def install():
    if not (host.data.blocky_enabled or host.data.k3s_master):
        return

    ##############
    # Keepalived #
    ##############

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

    if config.changed:
        systemd.daemon_reload(name="Reload systemd daemon")

    systemd.service(
        name="Restart and enable the keepalived service",
        service="keepalived.service",
        running=True,
        restarted=config.changed, # Trigger a restart only if the config changed
        enabled=True,
    )

    ############
    # Exporter #
    ############

    install_path = f"{host.data.keepalived_exporter_dir}/keepalived-exporter-{host.data.keepalived_exporter_version}"

    files.directory(
        name="Ensure exporter install dir exists",
        path=install_path,
        present=True,
        mode="755",
        recursive=True,
    )

    exporter_bin = f"{install_path}/keepalived-exporter"

    if host.get_fact(File, exporter_bin) is None:
        download_dest = f"{host.data.keepalived_exporter_tmp_dir}/keepalived-exporter.tar.gz"

        files.download(
            name="Download archive",
            src=exporter_download_url(),
            dest=download_dest,
            force=True, # always download the file, even if it already exists
        )

        server.shell(
            name="Expand archive",
            commands=[
                f"tar -C {install_path} -zxvf {download_dest}",
            ]
        )

    files.link(
        name=f"Link {host.data.keepalived_exporter_dir}/keepalived-exporter as {install_path}/keepalived-exporter",
        path=f"{host.data.keepalived_exporter_dir}/keepalived-exporter",
        target=f"{install_path}/keepalived-exporter",
    )

    exporter_unit = files.template(
        name="Configure keepalived-exporter service",
        src=resource_path("templates/keepalived/keepalived-exporter.service.j2"),
        dest="/etc/systemd/system/keepalived-exporter.service",
        mode="644",
        user="root",
        group="root",
        keepalived_exporter_dir=host.data.keepalived_exporter_dir,
    )

    if exporter_unit.changed:
        systemd.daemon_reload(name="Reload systemd daemon")

    systemd.service(
        name="Restart and enable the keepalived-exporter service",
        service="keepalived-exporter.service",
        running=True,
        restarted=exporter_unit.changed, # Trigger a restart only if the config or unit changed
        enabled=True,
    )

def exporter_download_url():
    arch = host.get_fact(Arch)
    if arch == "x86_64":
        arch = "amd64"

    return "{0}/v{1}/keepalived-exporter_{1}_{2}_{3}.tar.gz".format(
        host.data.keepalived_exporter_base_url,
        host.data.keepalived_exporter_version,
        host.get_fact(Os).lower(),
        arch,
    )
