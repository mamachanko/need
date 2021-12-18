SHELL := bash
.SHELLFLAGS := -euo pipefail -c
.ONESHELL :=
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules


.PHONY: test
test:
	# expected to succeed
	go run main.go --file test/satisfiable.yaml

	# expected to fail
	! go run main.go --file test/lacking.yaml

	# test output
	diff \
	  <(go run main.go --file test/lacking.yaml --file test/satisfiable.yaml) \
	  <(cat test/expected_output)

	# test output when failing fast
	diff \
	  <(go run main.go --file test/lacking.yaml --file test/satisfiable.yaml --fail-fast) \
	  <(cat test/expected_output_fail_fast)

.PHONY: example
example:
	@go run main.go \
	  --file example.yaml \
	  --fail-fast i
