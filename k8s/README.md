# Kubernetes

Installed using [K3S](https://docs.k3s.io/)

## Server Setup (carrot)

Setup k3s with default options:

```shell
curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="--tls-san carrot.lab --disable traefik" sh -
```

`INSTALL_K3S_EXEC="--tls-san carrot.lab"` ensures that the certificate generated for the API server is valid for `carrot.lab`.

A kubeconfig is written to `/etc/rancher/k3s/k3s.yaml`.


## Node Setup (potato)

Run on potato:

```shell
curl -sfL https://get.k3s.io | K3S_URL=https://192.168.1.42:6443 K3S_TOKEN=mynodetoken sh -
```

Note: `K3S_TOKEN` is stored in `/var/lib/rancher/k3s/server/node-token` on the server node.

## Argo CD

### Installation

https://argo-cd.readthedocs.io/en/stable/getting_started/

```shell
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

### LDAP integration

```shell
kubectl edit configmap argocd-cm -n argocd
```

Connectors are managed by [dex](https://dexidp.io/docs/connectors/ldap/)

The configuration is automatically reloaded after editing this file:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-cm
  namespace: argocd
  labels:
    app.kubernetes.io/name: argocd-cm
    app.kubernetes.io/part-of: argocd
data:
  url: https://argocd.kevingomez.fr

  dex.config: |
    connectors:
      - type: ldap
        id: homelab-ldap
        name: LDAP Homelab
        config:
          host: 'lldap.lldap.svc.cluster.local:3890'
          insecureNoSSL: true
          bindDN: "uid=argocd,ou=people,dc=home,dc=lab"
          bindPW: "$ldap.ldap_password"

          # Ldap user search attributes
          userSearch:
            baseDN: "ou=people,dc=home,dc=lab"
            filter: ""
            username: user_id
            idAttr: uuid
            emailAttr: mail
            nameAttr: display_name

          # Ldap group search attributes
          groupSearch:
            baseDN: "ou=groups,dc=home,dc=lab"
            filter: "(objectClass=group)"
            userMatchers:
            - userAttr: uuid
              groupAttr: member
            nameAttr: uuid
```

RBAC policy to give the admin role to members of the LDAP group "argocd_admin":

The configuration is automatically reloaded after editing this file:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-rbac-cm
  namespace: argocd
  labels:
    app.kubernetes.io/name: argocd-rbac-cm
    app.kubernetes.io/part-of: argocd
data:
  policy.csv: |
    g, argocd_admin, role:admin
  policy.default: role:readonly
  scopes: '[groups, email]'
```

### Accessing the UI

```shell
# Get the password for the "admin" account
k8sec list -n argocd argocd-initial-admin-secret
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

### Cluster bootstrapping

The "app-of-apps" pattern is used to simplify bootstrapping: an app is the entrypoint for all others.

This app can be added via the CLI or via the UI.

```shell
argocd app create apps \
    --dest-namespace argocd \
    --dest-server https://kubernetes.default.svc \
    --repo https://github.com/K-Phoen/homelab.git \
    --path k8s/apps
argocd app sync apps 
```

## 1password integration

* https://external-secrets.io/main/provider/1password-automation/
* https://github.com/1Password/connect/blob/a0a5f3d92e68497098d9314721335a7bb68a3b2d/README.md#quick-start

```shell
op connect server create homelab --vaults Homelab
op vault list
op connect token create homelab_k8s --server homelab --vaults o327uryxcflmqtawyqbm3ezzoq,r

kubectl create secret generic -n 1password op-credentials --from-literal=1password-credentials.json=$(cat 1password-credentials.json|base64 -w0)
kubectl create secret generic -n 1password onepassword-token --from-literal=token=token
```

Creating a secret: https://developer.1password.com/docs/k8s/k8s-operator/#kubernetes-secret-from-item

```yaml
apiVersion: onepassword.com/v1
kind: OnePasswordItem
metadata:
  name: secret-name
  namespace: foo
spec:
  itemPath: "vaults/Homelab/items/[ITEM]"
```

Where `[ITEM]` is the ID or title of the 1password item to synchronize.

## Monitoring

Create a few secrets:

```shell
kubectl create namespace monitoring
kubectl create secret generic -n monitoring grafana-k8s-monitoring \
    --from-literal=gcloud_metrics_user='ID' \
    --from-literal=gcloud_metrics_password='TOKEN' \
    --from-literal=gcloud_logs_user='ID' \
    --from-literal=gcloud_logs_password='TOKEN' \
    --from-literal=gcloud_otlp_user='ID' \
    --from-literal=gcloud_otlp_password='TOKEN' \
    --from-literal=gcloud_remote_fleet_management_user='ID' \
    --from-literal=gcloud_remote_fleet_management_password='TOKEN'
```

## Democratic CSI

Useful resources:

* https://github.com/democratic-csi/democratic-csi
* https://github.com/democratic-csi/democratic-csi/blob/master/examples/freenas-api-nfs.yaml
* https://jonathangazeley.com/2021/01/05/using-truenas-to-provide-persistent-storage-for-kubernetes/

Pre-requisites:

* create a `k8s` dataset in the `main` pool
* create a `k8s-nfs` user in TrueNAS
* `cp values-truenas-nfs-api{.redacted,}.yaml`
* define the API key to use in `values-truenas-nfs-api.yaml` (and tweak other values as needed)

```shell
cd democratic-csi
helm repo add democratic-csi https://democratic-csi.github.io/charts/
helm repo update
helm upgrade \
  --install \
  --create-namespace \
  --values values-truenas-nfs-api.yaml \
  --namespace democratic-csi \
  --set node.kubeletHostPath="/var/lib/kubelet" \
  truenas-nfs-api democratic-csi/democratic-csi
k get storageclasses # there should be a truenas-nfs-api-csi StorageClass
k apply -f test-pvc-nfs-api.yaml
k describe pvc -n democratic-csi test-claim-truenas-nfs-api
k delete -n democratic-csi test-claim-truenas-nfs-api
```