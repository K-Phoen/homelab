from pyinfra.context import config, host
from infra.tasks import alloy, blocky, keepalived, k3s, packages, tailscale, ufw

config.SUDO_PASSWORD = host.data.sudo_password

packages.install()
alloy.install()

if host.data.blocky_enabled:
    blocky.install()

if host.data.k3s_master:
    k3s.install()

keepalived.install()
tailscale.install()
ufw.install()
