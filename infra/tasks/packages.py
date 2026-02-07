from pyinfra.api.deploy import deploy
from pyinfra.operations import apt

@deploy("Install packages")
def install_packages():
    apt.packages(
        name="Install useful packages",
        packages=[
          "bat",
          "btop",
          "curl",
          "git",
          "lm-sensors",
          "vim",
          "dnsutils",
        ],
        no_recommends=True,
        update=True,
    )

    apt.packages(
        name="Install required packages",
        packages=[
          "firmware-linux-nonfree",
          "open-iscsi",
          "nfs-common",
        ],
        no_recommends=True,
        update=True,
    )
