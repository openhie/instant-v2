#!/bin/bash

readonly ACTION=$1
readonly MODE=$2

readonly STATEFUL_NODES=${STATEFUL_NODES:-"cluster"}

PACKAGE_PATH=$(
    cd "$(dirname "${BASH_SOURCE[0]}")" || exit
    pwd -P
)
readonly PACKAGE_PATH

ROOT_PATH="${PACKAGE_PATH}/.."
readonly ROOT_PATH

. "${ROOT_PATH}/utils/config-utils.sh"
. "${ROOT_PATH}/utils/docker-utils.sh"
. "${ROOT_PATH}/utils/log.sh"

main() {
    if [[ "${MODE}" == "dev" ]]; then
        log info "Running {{.Id}} in DEV mode"
        dev_compose_file="-c ${PACKAGE_PATH}/docker-compose.dev.yml"
    else
        log info "Running {{.Id}} in PROD mode"
        dev_compose_file=""
    fi

    if [[ "${ACTION}" == "init" ]] || [[ "${ACTION}" == "up" ]]; then
        config::set_config_digests "${PACKAGE_PATH}"/docker-compose.yml
        try "docker stack deploy -c ${PACKAGE_PATH}/docker-compose.yml $dev_compose_file {{.Stack}}" "Failed to deploy {{.Id}}"

        docker::await_container_startup {{.Id}}
        docker::await_container_status {{.Id}} Running

    elif [[ "${ACTION}" == "down" ]]; then
        try "docker service scale {{.Stack}}_{{.Id}}=0" "Failed to scale down {{.Id}}"

    elif [[ "${ACTION}" == "destroy" ]]; then
        docker::service_destroy {{.Id}}

    else
        log error "Valid options are: init, up, down, or destroy"
    fi
}

main "$@"
