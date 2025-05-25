#!/bin/bash

tag=local
repository=""

images=(
    "${repository}gateway:$tag" 
    "${repository}storage:$tag" 
    "${repository}storage-pg-migrate:$tag" 
)
paths=(
    "gateway/"
    "storage/"
    "storage/"
)
dockerfiles=(
    "Dockerfile"
    "Dockerfile"
    "migrate.Dockerfile"
)

build-debug-local() {
    echo "Bulding docker images..."
    for i in "${!images[@]}"; do
        if ! minikube image build "${paths[$i]}" \
            -f "${dockerfiles[$i]}" \
            -t "${images[$i]}"
        then
            echo "Building ${images[$i]} failed"
            exit 1
        fi
    done

    echo "Successfully built docker images!"
}

build-debug-local
