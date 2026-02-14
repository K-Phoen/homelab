from getpass import getpass
import privy

# Tools
decrypt_pwd = getpass("Password (decrypting secrets): ")
decrypt = lambda val: privy.peek(val, decrypt_pwd).decode("utf-8")

# Globals
node_exporter_enabled = True

# Alloy
gcloud_prometheus_username = "1021535"
gcloud_loki_username = "613373"
gcloud_token = decrypt("1$2$FS0wznSFBxAE07wR8CyVsxH0ev-p7BlyC1PGXZd8HNA=$Z0FBQUFBQnBpUVhJZGw2ZVlRQlJCX3RRODZ5LTVra3NNR0RISWg0OWgyby01RVQ1VTNfRFZFLUdaOEMxYjNOZWtLTnVfVFlFOVczV09xRkhPeFNNY283MVUzZkhveW1kYThoRmFNdFZqV2ZRc09tcWtVRnlGa0RwQ2FQNGlZTmxIZy1BQzBnRDVuUlZPejlETEdCUFN4ZGJRT05ESGlyaEdEbjJyNDd5cjBEZ3BfVXg0d25DNklaM1VTaHlicWZUTENOU09TNEMtOUltYzN5M29HMUxzc2tkYTBaZ3hoaWZ3RzhFSjdzWWljTTRISlhWN0VwQ2lwNGVkVmZmcndpTVVmX3dWc1NjTXFkQnU3UFFXWHFTWXlPVEF3MExDY3RWS3c9PQ==")

# Tailscale
tailscale_authkey = decrypt("1$2$_vg0h9wMwiK4p0kqF7P_Oiq0ifQR7Tlo6aAJzvzJ-Fg=$Z0FBQUFBQnBrT2JYRnE0ak42N3ZMNFpLMDFfY1JVM09iR09vcnNmbFlfREI3RjJ5YktoUGd2MkJZUWVfbzJsQzltYlR2TGZGTm9hOU9GQkZ1TVl5ZXhPTEN4MXFOMnFWM0tHT2xwR0NtYXN0ak9ERHdab1BtYTREMVpObEtBV2ZxLThXQU9PTjVfb05iSDlQLXFHVmhOeDk0dHdjR1plS2xnPT0=")
tailscale_args = [
    "--ssh",
    "--advertise-exit-node",
    "--advertise-tags=tag:homelab",
    "--auth-key=%s?ephemeral=false".format(tailscale_authkey),
    "--reset",
]