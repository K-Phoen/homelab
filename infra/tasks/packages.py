from pyinfra.api.deploy import deploy
from pyinfra.operations import apt

@deploy("Install useful and required packages")
def install():
    apt.packages(
        name="Install packages",
        packages=[
          # Useful
          "bat",
          "btop",
          "curl",
          "git",
          "lm-sensors",
          "vim",
          "dnsutils",
          # Required
          "firmware-linux-nonfree",
          "open-iscsi",
          "nfs-common",
        ],
        no_recommends=True,
        update=True,
        cache_time=3600, # seconds
    )
