SHELL := bash
.SHELLFLAGS := -euo pipefail -c
.ONESHELL :=
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

.DEFAULT_GOAL := test

.PHONY: test
test:
	DEV_DEPENDENCIES_CONFIG="./test/good.yml" ./devDependencies.sh
	! { DEV_DEPENDENCIES_CONFIG="./test/bad.yml" ./devDependencies.sh; }
