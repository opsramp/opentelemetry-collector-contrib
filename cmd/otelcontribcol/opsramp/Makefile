IMAGE    ?= varmasai10/otel
TAG      ?= latest
PLATFORM ?= linux/amd64
REPO_ROOT := $(abspath $(CURDIR)/../../..)

.PHONY: all image push

all: push

# The Dockerfile compiles otelcontribcol from the current sources, so the image always carries a fresh binary.
image:
	docker build --platform $(PLATFORM) \
		-f $(REPO_ROOT)/cmd/otelcontribcol/opsramp/Dockerfile \
		-t $(IMAGE):$(TAG) $(REPO_ROOT)

push: image
	docker push $(IMAGE):$(TAG)
