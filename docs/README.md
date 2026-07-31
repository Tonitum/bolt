# bolt

## Local Kubernetes

Requirements: Docker, `kind`, `kubectl`, and Helm.

```sh
brew install kind kubectl helm
kind create cluster --name bolt --config kind.yaml
```

Ensure `blt` resolves to `127.0.0.1` in `/etc/hosts`.

Build and load the application image:

```sh
docker build -t bolt:dev .
kind load docker-image bolt:dev --name bolt
```

Install Bolt:

```sh
helm upgrade --install bolt charts/bolt \
  --set image.repository=bolt \
  --set image.tag=dev \
  --set image.pullPolicy=IfNotPresent \
  --set service.type=NodePort \
  --set service.nodePort=30080
```

Check the resources:

```sh
kubectl get pods,pvc,service
```

Test the route:

```sh
curl http://blt/list
```
