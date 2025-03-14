# Kubernetes

Installed using [K3S](https://docs.k3s.io/)

## Server Setup (carrot)

Setup k3s with default options:

```shell
curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="--tls-san carrot.lab" sh -
```

`INSTALL_K3S_EXEC="--tls-san carrot.lab"` ensures that the certificate generated for the API server is valid for `carrot.lab`.

A kubeconfig is written to `/etc/rancher/k3s/k3s.yaml`.

## Argo CD

### Installation

https://argo-cd.readthedocs.io/en/stable/getting_started/

```shell
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
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
