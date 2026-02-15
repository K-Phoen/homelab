from pyinfra import config
from pyinfra.context import host
from infra.tasks import alloy, packages, tailscale

config.SUDO_PASSWORD = host.data.sudo_password

packages.install()
alloy.install()
tailscale.install()
