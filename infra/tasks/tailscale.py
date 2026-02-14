from pyinfra.context import host
from pyinfra.api.deploy import deploy
from pyinfra.operations import apt, files, server, systemd

from .defaults import DEFAULTS
from .utils import resource_path

@deploy("Install Tailscale", data_defaults=DEFAULTS)
def install():
    files.directory(
        name="Ensure keyrings directory exists",
        path="/usr/share/keyrings",
        user="root",
        group="root",
        mode="755",
        present=True,
    )

    files.download(
        name="Download GPG key",
        src="https://pkgs.tailscale.com/stable/debian/trixie.noarmor.gpg",
        dest="/usr/share/keyrings/tailscale-archive-keyring.gpg",
        user="root",
        group="root",
        mode="644",
    )

    files.put(
        name="Ensure Tailscale source file exists",
        src=resource_path("files/tailscale.sources"),
        dest="/etc/apt/sources.list.d/tailscale.sources",
        mode="644",
        user="root",
        group="root",
    )

    apt.packages(
        name="Install Tailscale package",
        packages=["tailscale"],
        no_recommends=True,
        update=True,
        cache_time=3600, # seconds
    )

    sysctl_ts = files.put(
        name="Ensure IP forwarding is enabled",
        src=resource_path("files/99-tailscale.conf"),
        dest="/etc/sysctl.d/99-tailscale.conf",
        mode="644",
        user="root",
        group="root",
    )

    if sysctl_ts.changed:
        server.shell(
            name="Load 99-tailscale.conf",
            commands=[
                "sysctl -p /etc/sysctl.d/99-tailscale.conf",
            ],
            _sudo=True,
        )

    server.shell(
        name="Run tailscale up",
        commands=[
            "tailscale up %s" % " ".join(host.data.tailscale_args),
        ],
        _sudo=True,
    )

    systemd.service(
        name="Enable and restart the Tailscaled service",
        service="tailscaled.service",
        running=True,
        enabled=True,
    )
