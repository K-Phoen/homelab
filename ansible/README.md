# Ansible

## Encrypting secrets

See `secrets.example` file for a reference of what's expected in the `secrets.encrypted` file.

```shell
cp secrets.example secrets.encrypted
ansible-vault encrypt secrets.encrypted
ansible-vault edit secrets.encrypted
```