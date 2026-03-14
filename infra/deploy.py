from pyinfra.context import config, host
from infra.tasks import alloy, blocky, packages, tailscale

config.SUDO_PASSWORD = host.data.sudo_password

packages.install()
alloy.install()
if host.data.blocky_enabled:
    blocky.install()
tailscale.install()
