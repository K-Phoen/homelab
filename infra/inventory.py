sandbox_hosts = [
    ("potato", {
        "ssh_hostname": "10.10.10.122",
        "ssh_port": 4242,
        "ssh_user": "kevin",
        "_sudo": True,

        "node_exporter_enabled": True,
    }),
]

homelab_hosts = [
    ("carrot", {
        "ssh_hostname": "10.10.10.64",
        "ssh_port": 4242,
        "ssh_user": "kevin",
        "_sudo": True,

        "blocky_enabled": True,
        "k3s_master": True,

        "keepalived_vrrp_priority_blocky": 100,
        "keepalived_vrrp_priority_k3s": 150, # higher is better
   }),
   ("bean", {
        "ssh_hostname": "10.10.10.91",
        "ssh_port": 4242,
        "ssh_user": "kevin",
        "_sudo": True,

        "blocky_enabled": True,
        "k3s_master": True,

        "keepalived_vrrp_priority_blocky": 150, # higher is better
        "keepalived_vrrp_priority_k3s": 100,
   }),
   ("celery", {
        "ssh_hostname": "10.10.10.14",
        "ssh_port": 4242,
        "ssh_user": "kevin",
        "_sudo": True,

        "blocky_enabled": True,
        "k3s_master": True,

        "keepalived_vrrp_priority_blocky": 80,
        "keepalived_vrrp_priority_k3s": 80,
   }),
]
