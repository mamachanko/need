#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"

: "${DEV_DEPENDENCIES_CONFIG:=./devDependencies.yml}"

DEV_DEPENDENCY_CHECKS=
EXIT_CODE=0

main() {
  read_dependency_checks
  check_and_print_status
  print_conditions_and_exit
}

read_dependency_checks() {
  readarray DEV_DEPENDENCY_CHECKS < <(
    yq eval \
      --output-format=json \
      --indent=0 \
      .dependencies[] \
      "$DEV_DEPENDENCIES_CONFIG"
  )
}

check_and_print_status() {
  cat <<EOF
$(cat $DEV_DEPENDENCIES_CONFIG)
status:
  dependencies:
EOF

  for check in "${DEV_DEPENDENCY_CHECKS[@]}"; do
    name="$(echo "$check" | jq -r .name)"
    help="$(echo "$check" | jq -r .help)"
    testCmd="$(echo "$check" | jq -r .testCmd)"

    if output="$(bash -euo pipefail -c "$testCmd" 2>&1)"; then
      cat <<EOF
  - name: $name
    status: ✅
EOF
    else
      cat <<EOF
  - name: $name
    testCmd: |
$(echo "$testCmd" | sed 's/^/      /g')
    output: |
$(echo "$output" | sed 's/^/      /g')
    status: ❌
    help: |
      $help
EOF
      EXIT_CODE=1
    fi
  done
}

print_conditions_and_exit() {
  if [ $EXIT_CODE = 0 ]; then
    cat <<EOF
  conditions:
    - type: AllDependenciesSatisfied
      lastCheckTime: $(date)
      status: ✅
EOF
  else
    cat <<EOF
  conditions:
    - type: AllDependenciesSatisfied
      lastCheckTime: $(date)
      status: ❌
      message: see {status.dependencies[].help} for help
EOF
  fi
  exit $EXIT_CODE
}

main
