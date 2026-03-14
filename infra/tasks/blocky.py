from pyinfra.context import host
from pyinfra.api.deploy import deploy
from pyinfra.facts.server import Arch, Os
from pyinfra.facts.files import File
from pyinfra.operations import files, server, systemd

from .defaults import DEFAULTS
from .utils import resource_path

@deploy("Install Blocky", data_defaults=DEFAULTS)
def install():
    server.group(
        name="Create blocky group",
        group="blocky",
        system=True,
    )

    server.user(
        name="Create blocky user",
        user="blocky",
        group="blocky",
        create_home=False,
        system=True,
    )

    files.directory(
        name="Ensure install dir exists",
        path=host.data.blocky_install_path,
        present=True,
        user="blocky",
        group="blocky",
        mode="755",
        recursive=True,
    )

    blocky_bin = f"{host.data.blocky_install_path}/blocky"

    if host.get_fact(File, blocky_bin) is None:
        download_dest = f"{host.data.blocky_tmp_dir}/blocky.tar.gz"

        files.download(
            name="Download blocky archive",
            src=blocky_download_url(),
            dest=download_dest,
            user="blocky",
            group="blocky",
            # TODO: parse checksums file and verify checksum
            #sha256sum=urllib.request.urlopen("{}/{}/blocky_checksums.txt".format(host.data.blocky_base_url, host.data.blocky_version)).read(),
        )

        server.shell(
            name="Expand blocky archive",
            commands=[
                f"tar -C {host.data.blocky_install_path} -zxvf {download_dest}",
            ]
        )

        files.file(
            name="Ensure owner and mode of blocky binary",
            path=f"{host.data.host.data.blocky_install_path}/blocky",
            user="blocky",
            group="blocky",
            mode="755",
        )

    files.link(
        name=f"Link {host.data.blocky_dir}/blocky as {host.data.blocky_install_path}/blocky",
        path=f"{host.data.blocky_dir}/blocky",
        target=f"{host.data.blocky_install_path}/blocky",
        user="blocky",
        group="blocky",
    )

    blocky_config = files.put(
        name="Configure blocky",
        src=resource_path("files/blocky/config.yaml"),
        dest=f"{host.data.blocky_dir}/config.yaml",
        mode="644",
        user="blocky",
        group="blocky",
    )

    blocky_unit = files.template(
        name="Configure blocky service",
        src=resource_path("templates/blocky/blocky.service.j2"),
        dest="/etc/systemd/system/blocky.service",
        mode="644",
        user="root",
        group="root",
    )

    systemd.service(
        name="Restart and enable the blocky service",
        service="blocky.service",
        running=True,
        restarted=blocky_config.changed or blocky_unit.changed, # Trigger a restart only if the config or unit changed
        enabled=True,
    )

def blocky_download_url():
    arch = host.get_fact(Arch)
    if arch == "aarch64":
        arch = "arm64"

    return "{0}/{1}/blocky_{1}_{2}_{3}.tar.gz".format(
        host.data.blocky_base_url,
        host.data.blocky_version,
        host.get_fact(Os),
        arch,
    )
