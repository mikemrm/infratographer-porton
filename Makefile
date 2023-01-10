# Runtime
LICENSE_PATH?=$(PWD)/krakend-license

# Container variables
IMAGE_REGISTRY?=ghcr.io
IMAGE_REF?=$(IMAGE_REGISTRY)/infratographer/porton/porton
IMAGE_TAG?=latest
IMAGE?=$(IMAGE_REF):$(IMAGE_TAG)

CNT_CMD?=docker
CNT_RUN_CMD?=$(CNT_CMD) run --rm -it
CNT_BUILD_CMD?=docker buildx build

# Helpers
RUN_BUILD?=true

# Targets
.PHONY: image
image:
ifeq ($(RUN_BUILD),true)
	@echo "Building Portón container image"
	$(CNT_BUILD_CMD) -t $(IMAGE) .
else
	@echo "Skipping Portón container image build"
endif

.PHONY: klicense
klicense:
	@echo "Ensuring krakend license"
	test -f $(LICENSE_PATH) || (echo "Missing krakend license file" && exit 1)

.PHONY: run
run: image klicense
	@echo "Running Portón API gateway"
	$(CNT_RUN_CMD) -p 8080:8080 \
		-v $(LICENSE_PATH):/etc/krakend/LICENSE \
		-v $(PWD)/tests/data/krakend-minimal-config.json:/etc/krakend/krakend.json \
		$(IMAGE) run --config /etc/krakend/krakend.json