app := krawler
version := 0.1.0

oci_image := quay.io/maxgio92/$(app)

bins := go golangci-lint gofumpt
commands := version list

define declare_binpaths
$(1) = $(shell command -v 2>/dev/null $(1))
endef

.PHONY: build
build:
	@$(go) build .

.PHONY: test
test:
	@go test -v -cover -gcflags=-l ./...

.PHONY: lint
lint: golangci-lint
	@$(golangci-lint) run ./...

.PHONY: golangci-lint
golangci-lint:
	@$(go) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.0

.PHONY: gofumpt
gofumpt:
	@$(go) install mvdan.cc/gofumpt@v0.3.1

$(foreach bin,$(bins),\
	$(eval $(call declare_binpaths,$(bin)))\
)

.PHONY: oci/build
oci/build:
	@docker build . -t $(oci_image):$(version) -f Containerfile

.PHONY: oci/push
oci/push: oci/build
	@docker push $(oci_image):$(version)

.PHONY: clean
clean:
	@rm -f $(app)

.PHONY: help
help: list

.PHONY: list
list:
	@LC_ALL=C $(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'
