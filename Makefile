app := krawler
version := 0.1.0

oci_image := quay.io/maxgio92/$(app)

bins := go golangci-lint gofumpt aws

DISTROS ?= amazonlinux amazonlinux2 amazonlinux2022 centos debian ubuntu

RESULTS_DIR := e2e/results

BUCKET_NAME := krawler-kernel-releases

define declare_binpaths
$(1) = $(shell command -v 2>/dev/null $(1))
endef

define gen_run_targets
.PHONY: run/$(1)
run/$(1):
	@rm -rf $(RESULTS_DIR)/$(1) 2>/dev/null || true
	@mkdir -p $(RESULTS_DIR)/$(1)

	@echo -n "$(1) with default configuration: "
	@./$(app) list $(1) \
	  -o json \
	  > $(RESULTS_DIR)/$(1)/index.json 2> $(RESULTS_DIR)/$(1)/krawler.log
	@echo "$$$$(jq length $(RESULTS_DIR)/$(1)/index.json) releases found."
endef

$(foreach bin,$(bins),\
	$(eval $(call declare_binpaths,$(bin)))\
)

define gen_e2e_targets
.PHONY: e2e/$(1)
e2e/$(1): run/$(1)
	@echo -n "$(1) with custom configuration (full): "
	@./$(app) list $(1) \
	  -c testdata/$(1).yaml \
	  -o json \
	  > $(RESULTS_DIR)/$(1)/index_custom.json 2> $(RESULTS_DIR)/$(1)/krawler_custom.log
	@echo "$$$$(jq length $(RESULTS_DIR)/$(1)/index_custom.json) releases found."

	@echo -n "$(1) with custom configuration (no repositories): "
	@./$(app) list $(1) \
	  -c testdata/$(1)-norepos.yaml \
	  -o json \
	  > $(RESULTS_DIR)/$(1)/index_custom_norepos.json 2> $(RESULTS_DIR)/$(1)/krawler_custom_norepos.log
	@echo "$$$$(jq length $(RESULTS_DIR)/$(1)/index_custom_norepos.json) releases found."

	@{ DEFAULT=$$$$(jq length $(RESULTS_DIR)/$(1)/index.json) \
	CUSTOM=$$$$(jq length $(RESULTS_DIR)/$(1)/index_custom.json) \
	CUSTOM_NOREPOS=$$$$(jq length $(RESULTS_DIR)/$(1)/index_custom_norepos.json); \
		[[ $$$$DEFAULT == $$$$CUSTOM ]] && \
		[[ $$$$CUSTOM == $$$$CUSTOM_NOREPOS ]] && \
		echo "$(1) OK"; \
	} \
	|| { echo "$(1) KO"; exit 1; }
endef

define gen_publish_targets
.PHONY: publish/$(1)
publish/$(1): run/$(1)
	$(aws) s3 sync $(RESULTS_DIR)/$(1)/ s3://$(BUCKET_NAME)/$(1)/
endef

$(foreach distro,$(DISTROS),\
	$(eval $(call gen_run_targets,$(distro)))\
	$(eval $(call gen_e2e_targets,$(distro)))\
	$(eval $(call gen_publish_targets,$(distro)))\
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
	@docker build . -t $(oci_image):$(version) -f Dockerfile

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

.PHONY: run
run: clean build $(patsubst %,run/%,$(DISTROS))

.PHONY: e2e
e2e: clean build $(patsubst %,e2e/%,$(DISTROS))

.PHONY: publish
publish: clean build $(patsubst %,publish/%,$(DISTROS))

PACKAGE_NAME          := github.com/maxgio92/$(app)
GOLANG_CROSS_VERSION  ?= v$(shell sed -nE 's/go[[:space:]]+([[:digit:]]\.[[:digit:]]+)/\1/p' go.mod)

.PHONY: release
release:
	@if [ ! -f ".release-env" ]; then \
		echo "\033[91m.release-env is required for release\033[0m";\
		exit 1;\
	fi
	docker run \
		--rm \
		-e CGO_ENABLED=1 \
		--env-file .release-env \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --rm-dist
