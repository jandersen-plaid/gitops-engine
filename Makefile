# docker image publishing options
DOCKER_PUSH?=false
IMAGE_NAMESPACE?=
IMAGE_TAG?=latest
ifdef IMAGE_NAMESPACE
IMAGE_PREFIX=${IMAGE_NAMESPACE}/
endif

CLIENT_K8S_VERSION = $(shell grep "client-go" go.mod | head -n 1 | cut -d ' ' -f2)

.PHONY: generate
generate: agent-manifests

.PHONY: test
test:
	go test -race ./... -coverprofile=coverage.out

.PHONY: lint
lint:
	golangci-lint run

.PHONY: agent-image
agent-image:
	docker build -t $(IMAGE_PREFIX)gitops-agent . -f Dockerfile
	@if [ "$(DOCKER_PUSH)" = "true" ] ; then docker push $(IMAGE_PREFIX)gitops-agent:$(IMAGE_TAG) ; fi

.PHONY: agent-manifests
agent-manifests:
	kustomize build ./agent/manifests/cluster-install > ./agent/manifests/install.yaml
	kustomize build ./agent/manifests/namespace-install > ./agent/manifests/install-namespaced.yaml

.PHONY: generate-mocks
generate-mocks:
	go generate -x -v "github.com/argoproj/gitops-engine/pkg/utils/tracing/tracer_testing"

.PHONY: update-parser
update-parser:
	curl -L https://raw.githubusercontent.com/kubernetes/client-go/$(CLIENT_K8S_VERSION)/applyconfigurations/internal/internal.go > pkg/utils/kube/scheme/internal/parser.go
