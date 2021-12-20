SHELL := bash
.SHELLFLAGS := -euo pipefail -c
.ONESHELL :=
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

GINKGO := go run github.com/onsi/ginkgo/ginkgo
MAIN := go run cmd/main.go

.PHONY: test
test: unit-test integration-test

.PHONY: unit-test
unit-test:
	$(GINKGO) \
	  -r \
	  --randomizeAllSpecs \
	  --randomizeSuites \
	  --failOnPending \
	  --cover \
	  --trace \
	  --progress

# TODO migrate to ginkgo
.PHONY: integration-test
integration-test:
	# expected to succeed
	$(MAIN) --file test/satisfiable.yaml

	# expected to fail
	! $(MAIN) --file test/lacking.yaml

	# test output
	diff \
	  <($(MAIN) --file test/lacking.yaml --file test/satisfiable.yaml) \
	  <(cat test/expected_output)

	# read from stdin
	! cat test/lacking.yaml | $(MAIN) --file -
	diff \
	  <($(MAIN) --file test/satisfiable.yaml) \
	  <(cat test/satisfiable.yaml | $(MAIN) --file -)

.PHONY: example
example:
	@$(MAIN) --file example.yaml
