build:
    @./build-docker.sh

run: build
    @echo "Installing local Helm chart..."
    @if [[ $(kubectl config current-context) != 'minikube' ]]; then \
        echo "ERROR: You should set minikube context: 'kubectl config use-context minikube'"; \
        exit 1; \
    fi
    helm install cw2 ./helm --values helm/values-local.yaml
    minikube tunnel

stop:
    helm uninstall cw2
