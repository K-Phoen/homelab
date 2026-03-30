from pyinfra.context import config, host
from pyinfra.operations import apt

config.SUDO_PASSWORD = host.data.sudo_password

apt.update(
    name="apt update",
    cache_time=3600,
)

apt.upgrade(
    name="apt upgrade",
    auto_remove=True,
)

apt.dist_upgrade(
    name="apt dist-upgrade",
    auto_remove=True,
)
