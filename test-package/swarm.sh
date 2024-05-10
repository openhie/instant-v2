#!/bin/bash

declare ACTION=""
declare MODE=""
declare COMPOSE_FILE_PATH=""
declare UTILS_PATH="/home/ryanc/git/instant-v2/utils"
declare SERVICE_NAMES=()

function init_vars() {
  ACTION=$1
  MODE=$2

  COMPOSE_FILE_PATH=$(
    cd "$(dirname "${BASH_SOURCE[0]}")" || exit
    pwd -P
  )

  # UTILS_PATH="${COMPOSE_FILE_PATH}/../utils"

  SERVICE_NAMES=(
    "test-package"
  )

  readonly ACTION
  readonly MODE
  readonly COMPOSE_FILE_PATH
  readonly UTILS_PATH
  readonly SERVICE_NAMES
}

# shellcheck disable=SC1091
function import_sources() {
  source "${UTILS_PATH}/docker-utils.sh"
  source "${UTILS_PATH}/log.sh"
}

main() {
  init_vars "$@"
  import_sources

  log debug "Debug message"
  log info "Info message"
  log warn "Warning message"
  log error "Error message"

  try 'false' 'catch' 'Failed to run false (expected)'

  if [[ "${ACTION}" == "init" ]] || [[ "${ACTION}" == "up" ]]; then
    log info "Running package in Single node mode"

  elif [[ "${ACTION}" == "down" ]]; then
    log info "Scaling down package"

    docker::scale_services_down "${SERVICE_NAMES[@]}"
  elif [[ "${ACTION}" == "destroy" ]]; then
    log info "Destroying package"

  else
    log error "Valid options are: init, up, down, or destroy"
  fi
}

main "$@"
