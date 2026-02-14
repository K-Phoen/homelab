from infra.tasks import alloy, packages, tailscale

packages.install()
alloy.install()
tailscale.install()
