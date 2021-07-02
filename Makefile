NAME=k8sbootstrap
BINARY=terraform-provider-${NAME}

.PHONY: docs

build:
	go build -o ${BINARY}

docs:
	go generate

test:
	go test ./... -v -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test ./... -v -timeout 120m
