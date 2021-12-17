#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"

: "${DEV_DEPENDENCIES_CONFIG:=./devDependencies.yml}"

DEV_DEPENDENCIES=
EXIT_CODE=0

main() {
  case "${1:-"check"}" in
  check)
    read_dev_dependencies
    check_dev_dependencies
    print_conditions_and_exit
    ;;
  install)
    read_dev_dependencies
    reconcile_dev_dependencies
    print_conditions_and_exit
    ;;
  usage)
    echo "usage: $0 [check (default)|install|usage]"
    exit
    ;;
  esac

}

read_dev_dependencies() {
  readarray DEV_DEPENDENCIES < <(
    yq eval \
      --output-format=json \
      --indent=0 \
      .dependencies[] \
      "$DEV_DEPENDENCIES_CONFIG"
  )
}

check_dev_dependencies() {
  cat <<EOF
$(cat $DEV_DEPENDENCIES_CONFIG)
status:
  dependencies:
EOF

  for check in "${DEV_DEPENDENCIES[@]}"; do
    name="$(echo "$check" | jq -r .name)"
    help="$(echo "$check" | jq -r .help)"
    testCmd="$(echo "$check" | jq -r .testCmd)"

    if testCmdOutput="$(bash -euo pipefail -c "$testCmd" 2>&1)"; then
      cat <<EOF
  - name: $name
    status: ✅
EOF
    else
      cat <<EOF
  - name: $name
    testCmd: |
$(echo "$testCmd" | sed 's/^/      /g')
    testCmdOutput: |
$(echo "$testCmdOutput" | sed 's/^/      /g')
    status: ❌
    help: |
      $help
EOF
      EXIT_CODE=1
    fi
  done
}

reconcile_dev_dependencies() {
  cat <<EOF
$(cat $DEV_DEPENDENCIES_CONFIG)
status:
  dependencies:
EOF

  for check in "${DEV_DEPENDENCIES[@]}"; do
    name="$(echo "$check" | jq -r .name)"
    help="$(echo "$check" | jq -r .help)"
    testCmd="$(echo "$check" | jq -r .testCmd)"
    installCmd="$(echo "$check" | jq -r .installCmd)"

    if {
      installCmdOutput="$(bash -euo pipefail -c "$installCmd" 2>&1)" &&
        testCmdOutput="$(bash -euo pipefail -c "$testCmd" 2>&1)"
    }; then
      cat <<EOF
  - name: $name
    status: ✅
EOF
    else
      cat <<EOF
  - name: $name
    installCmd: |
$(echo "$installCmd" | sed 's/^/      /g')
    installCmdOutput: |
$(echo "$installCmdOutput" | sed 's/^/      /g')
    testCmd: |
$(echo "$testCmd" | sed 's/^/      /g')
    testCmdOutput: |
$(echo "$testCmdOutput" | sed 's/^/      /g')
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

main "$@"
