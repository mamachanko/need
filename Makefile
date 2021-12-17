SHELL := bash
.SHELLFLAGS := -euo pipefail -c
.ONESHELL :=
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

.DEFAULT_GOAL := test

.PHONY: test
test: test-check test-install

.PHONY: test-check
test-check:
	DEV_DEPENDENCIES_CONFIG="./test/good-check.yml" ./devDependencies.sh
	! { DEV_DEPENDENCIES_CONFIG="./test/bad-check.yml" ./devDependencies.sh; }

.PHONY: test-install
test-install:
	DEV_DEPENDENCIES_CONFIG="./test/good-install.yml" ./devDependencies.sh install
	DEV_DEPENDENCIES_CONFIG="./test/good-install.yml" ./devDependencies.sh
	! { DEV_DEPENDENCIES_CONFIG="./test/bad-install.yml" ./devDependencies.sh install; }
