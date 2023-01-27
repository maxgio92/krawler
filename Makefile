app := krawler
version := 0.1.0

oci_image := quay.io/maxgio92/$(app)

bins := go golangci-lint gofumpt
commands := version list

DISTROS ?= amazonlinux amazonlinux2 amazonlinux2022 centos debian

define declare_binpaths
$(1) = $(shell command -v 2>/dev/null $(1))
endef

define gen_e2e_targets
e2e/$(1):
	@rm -f e2e/results/$(1).* 2>/dev/null || true
	@rm -f e2e/results/$(1)_custom.* 2>/dev/null || true
	@./$(app) list $(1) \
	  -o json \
	  > e2e/results/$(1).json 2> e2e/results/$(1).log
	@echo "$(1) with default configuration: $$$$(jq length e2e/results/$(1).json) releases found."
	@./$(app) list $(1) \
	  -c testdata/$(1).yaml \
	  -o json \
	  > e2e/results/$(1)_custom.json 2> e2e/results/$(1)_custom.log
	@echo "$(1) with custom configuration (full): $$$$(jq length e2e/results/$(1)_custom.json) releases found."
	@./$(app) list $(1) \
	  -c testdata/$(1)-norepos.yaml \
	  -o json \
	  > e2e/results/$(1)_custom_norepos.json 2> e2e/results/$(1)_custom_norepos.log
	@echo "$(1) with custom configuration (no repositories): $$$$(jq length e2e/results/$(1)_custom_norepos.json) releases found."
endef

$(foreach DISTRO,$(DISTROS),\
	$(eval $(call gen_e2e_targets,$(DISTRO)))\
)

$(foreach bin,$(bins),\
	$(eval $(call declare_binpaths,$(bin)))\
)

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

.PHONY: e2e
e2e: $(patsubst %,e2e/%,$(DISTROS))
