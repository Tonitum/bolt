# bolt

## Local Kubernetes

Requirements: Docker, `kind`, `kubectl`, and Helm.

```sh
brew install kind kubectl helm
kind create cluster --name bolt
```

Build and load the application image:

```sh
docker build -t bolt:dev .
kind load docker-image bolt:dev --name bolt
```

Install Envoy Gateway and Gateway API resources:

```sh
scripts/install-envoy-gateway.sh
helm upgrade --install test-infra charts/test-infra
```

Install Bolt:

```sh
helm upgrade --install bolt charts/bolt \
  --set image.repository=bolt \
  --set image.tag=dev \
  --set image.pullPolicy=IfNotPresent
```

Check the resources:

```sh
kubectl get pods,pvc,gateway,httproute
kubectl describe httproute bolt-bolt
```

Forward the Gateway service:

```sh
export ENVOY_SERVICE=$(kubectl get service \
  --namespace envoy-gateway-system \
  --selector=gateway.envoyproxy.io/owning-gateway-namespace=default,gateway.envoyproxy.io/owning-gateway-name=gateway \
  --output=jsonpath='{.items[0].metadata.name}')

kubectl port-forward --namespace envoy-gateway-system service/$ENVOY_SERVICE 8080:80
```

Test the route:

```sh
curl http://localhost:8080/list
```
