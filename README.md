# Designing Software, control work 2

## Functions

- upload file
- get file
- statistics of the file (A. service)
- 100% plagiat check (A. servcie)
- (optional) visualization

## Desicions explanation

- I decided to store `txt` files just in `TEXT` column in PostgreSQL because they are not meant to be big.
- I use CloudNativePG  operator because it makes using pg in Kubernetes much easier, I don't wanna complain about PVC, replication and related things.
- Minio is not persistent because it is not needed here

## Prerequests

Backend is cloud native and ready to be ran in Kubernetes cluster, so here we go.

Install `minikube`: [tap](https://minikube.sigs.k8s.io/docs/start/?arch=%2Fmacos%2Farm64%2Fstable%2Fbinary+download)

Install `helm`: [tap](https://helm.sh/docs/intro/install/)

Enable Ingress controller for minikube cluster

```bash
minikube addons enable ingress
```

Start local Kubernetes cluster

```bash
minikube start
```

Install CloudNativePG operator into local cluster

```bash
kubectl apply --server-side -f \
  https://raw.githubusercontent.com/cloudnative-pg/cloudnative-pg/release-1.25/releases/cnpg-1.25.1.yaml
```

## Run

Build docker images

```bash
./build-docker.sh
```

Install Helm shart

```bash
helm install cw2 ./helm --values helm/values-local.yaml
```

Make minikube tunnel gateway traffic

```bash
minikube tunnel
```

Check if it is working

```bash
curl http://0.0.0.0/ping
```
