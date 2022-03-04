bins := go cobra
commands := version list

AUTHOR := 
LICENSE := 

COBRA_FLAGS :=

.PHONY: build
build:
	@$(go) build .

.PHONY: init
init:
	@$(go) mod init krawler
	@$(cobra) init $(COBRA_FLAGS)

.PHONY: lint
lint:
	@golangci-lint run ./...

define declare_binpaths
$(1) := $(shell command -v 2>/dev/null $(1))
endef

define declare_command_build_targets
.PHONY: cmd/$(1)
cmd/$(1):
	@test -f cmd/$(1).go || $(cobra) add $(COBRA_FLAGS) $(1)
endef

$(foreach bin,$(bins),\
	$(eval $(call declare_binpaths,$(bin)))\
)

$(foreach command,$(commands),\
	$(eval $(call declare_command_build_targets,$(command)))\
)
