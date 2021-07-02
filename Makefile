HOSTNAME=github.com
NAMESPACE=nimbolus
NAME=k8sbootstrap
BINARY=terraform-provider-${NAME}
VERSION=0.1.0
OS_ARCH=linux_amd64

default: build

.PHONY: docs

build:
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

docs:
	go generate

test:
	go test ./... -v $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
