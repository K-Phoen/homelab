from pyinfra.context import host
from pyinfra.api.command import StringCommand
from pyinfra.api.deploy import deploy
from pyinfra.api.operation import operation
from pyinfra.operations import apt, files, server

from .defaults import DEFAULTS
from .utils import resource_path


@deploy("Setup UFW", data_defaults=DEFAULTS)
def install():
    apt.packages(
        name="Install ufw",
        packages=["ufw"],
        update=True,
    )

    files.put(
        name="Setup default values",
        src=resource_path("files/ufw/defaults"),
        dest="/etc/default/ufw",
        user="root",
        group="root",
        mode="644",
    )

    server.shell(name="Enable logging", commands=["ufw logging medium"])
    server.shell(name="Disable ufw", commands=["ufw --force disable"])
    server.shell(name="Reset rules", commands=["ufw --force reset"])

    # Unconditional rules
    ufw_allow(name="Allow all access to port 4242 - Custom SSH", port=4242, src="10.10.20.1/24")

    # Conditional blocky rules
    if host.data.blocky_enabled:
        ufw_allow(name="Allow keepalived multicast", src="10.10.10.1/24", dest="224.0.0.18")
        ufw_allow(name="Allow all access to port 53 - DNS TCP", port=53, proto="tcp")
        ufw_allow(name="Allow all access to port 53 - DNS UDP", port=53, proto="udp")

    # k3s rules
    if host.data.k3s_master:
        ufw_allow(name="Allow access to k8s apiserver", port=6443, proto="tcp", src="10.10.20.1/24")
        ufw_allow(name="Allow access to k8s apiserver - node to node", port=6443, proto="tcp", src="10.10.10.1/24")
        ufw_allow(name="k3s HA with embedded etcd - 2379", port=2379, proto="tcp", src="10.10.10.1/24")
        ufw_allow(name="k3s HA with embedded etcd - 2380", port=2380, proto="tcp", src="10.10.10.1/24")
        ufw_allow(name="Open pods networking", src="10.42.0.0/16")
        ufw_allow(name="Open services networking", src="10.43.0.0/16")
        ufw_allow(name="Open Flannel vxlan", port=8472, src="10.10.10.1/24")
        ufw_allow(name="Open Kubelet metrics", port=10250, src="10.10.10.1/24")
        ufw_allow(name="Open node-exporter metrics", port=9100, src="10.10.10.1/24")
        ufw_allow(name="Open metallb memberlist", port=7946, src="10.10.10.1/24")

    # k3s service rules
    if host.data.k3s_master: # TODO: use something else than `k3s_master`
        ufw_allow(name="Allow all access to port 80 - HTTP", port=80, proto="tcp")
        ufw_allow(name="Allow all access to port 443 - HTTPS", port=443, proto="tcp")
        ufw_allow(name="Allow all access to port 22 - gitea", port=22, proto="tcp")
        ufw_allow(name="Allow transmission port - TCP", port=51413, proto="tcp")
        ufw_allow(name="Allow transmission port - UDP", port=51413, proto="udp")

    server.shell(name="Enable UFW", commands=["ufw --force enable"])


@operation()
def ufw_allow(port: int | None = None, proto: str | None = None, src: str | None = None, dest: str | None = None):
    parts = ["ufw", "allow"]
    if src:
        parts += ["from", src]
    else:
        parts += ["from any"]

    if dest:
        parts += ["to", dest]
    else:
        parts += ["to any"]

    if port:
        parts += ["port", str(port)]

    if proto:
        parts += ["proto", proto]

    yield StringCommand(*parts)
