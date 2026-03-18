from getpass import getpass
import privy

# Tools
decrypt_pwd = getpass("Password (decrypting secrets): ")
decrypt = lambda val: privy.peek(val, decrypt_pwd).decode("utf-8")

# Globals
sudo_password = decrypt("1$2$1ugAONdjXv1or6XhYqJolbWbGYfovmxj8hJAz4h3wQ8=$Z0FBQUFBQnBramw5RmUwMUxRUHdxODFFSmlIOHlZRDZxQk4yUFBQY1dFZnlhX1loY0hiOEVhcXl6aVYzdnRUNXZVdUZ1Q19MRHhLcXRWelR6amd3Z3pKUTdIbUZRYUNTQ0E9PQ==")

# Alloy
node_exporter_enabled = False
gcloud_prometheus_username = "1021535"
gcloud_loki_username = "613373"
gcloud_token = decrypt("1$2$FS0wznSFBxAE07wR8CyVsxH0ev-p7BlyC1PGXZd8HNA=$Z0FBQUFBQnBpUVhJZGw2ZVlRQlJCX3RRODZ5LTVra3NNR0RISWg0OWgyby01RVQ1VTNfRFZFLUdaOEMxYjNOZWtLTnVfVFlFOVczV09xRkhPeFNNY283MVUzZkhveW1kYThoRmFNdFZqV2ZRc09tcWtVRnlGa0RwQ2FQNGlZTmxIZy1BQzBnRDVuUlZPejlETEdCUFN4ZGJRT05ESGlyaEdEbjJyNDd5cjBEZ3BfVXg0d25DNklaM1VTaHlicWZUTENOU09TNEMtOUltYzN5M29HMUxzc2tkYTBaZ3hoaWZ3RzhFSjdzWWljTTRISlhWN0VwQ2lwNGVkVmZmcndpTVVmX3dWc1NjTXFkQnU3UFFXWHFTWXlPVEF3MExDY3RWS3c9PQ==")

# Blocky
blocky_enabled = False
blocky_version = "v0.29.0"
blocky_dir = "/opt/blocky"
blocky_install_path = "{}/blocky-{}".format(blocky_dir, blocky_version)
blocky_base_url = "https://github.com/0xERR0R/blocky/releases/download"
blocky_tmp_dir = "/tmp"

# k3s

k3s_master = False

# Keepalived

keepalived_vip_blocky = '10.10.10.211'

# arbitrary unique number from 1 to 255
# used to differentiate multiple instances of vrrpd
keepalived_virtual_router_id_blocky = 42

# for electing MASTER, highest priority wins.
# The valid range of values for priority is [1-255], with priority
# 255 meaning "address owner".
# To be MASTER, it is recommended to make this 50 more than on
# other machines.
keepalived_vrrp_priority_blocky = 100

keepalived_vip_k3s = '10.10.10.212'
keepalived_vrrp_priority_k3s = 100
keepalived_virtual_router_id_k3s = 43

keepalived_exporter_version = "1.7.0"
keepalived_exporter_dir = "/opt/keepalived-exporter"
keepalived_exporter_install_path = f"{keepalived_exporter_dir}/keepalived-exporter-{keepalived_exporter_version}"
keepalived_exporter_base_url = "https://github.com/mehdy/keepalived-exporter/releases/download"
keepalived_exporter_tmp_dir = "/opt/keepalived-exporter"

# Tailscale
tailscale_authkey = decrypt("1$2$_vg0h9wMwiK4p0kqF7P_Oiq0ifQR7Tlo6aAJzvzJ-Fg=$Z0FBQUFBQnBrT2JYRnE0ak42N3ZMNFpLMDFfY1JVM09iR09vcnNmbFlfREI3RjJ5YktoUGd2MkJZUWVfbzJsQzltYlR2TGZGTm9hOU9GQkZ1TVl5ZXhPTEN4MXFOMnFWM0tHT2xwR0NtYXN0ak9ERHdab1BtYTREMVpObEtBV2ZxLThXQU9PTjVfb05iSDlQLXFHVmhOeDk0dHdjR1plS2xnPT0=")
tailscale_args = [
    "--ssh",
    "--advertise-exit-node",
    "--advertise-tags=tag:homelab",
    "--auth-key=%s?ephemeral=false".format(tailscale_authkey),
    "--reset",
]
