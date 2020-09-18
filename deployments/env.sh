#!/bin/bash

set -e

cd "$(dirname "$0")"

env=$1
cmd=$2
shift 2 || true

[ -f ${env}.env ] && source ${env}.env
[ -f ~/.tl/${env}.env ] && source ~/.tl/${env}.env

project_name=${PROJECT_NAME:-${env}}

compose="docker-compose -p ${project_name} -f docker-compose.yaml -f docker-compose.${env}.yaml"

case "${cmd}" in
up)
    services=$*
    eval "${compose} run db_migrations up"
    eval "${compose} up -d"
    ;;
down|ps|logs|exec|start|stop|restart|pull)
    args=$*
    eval "${compose} ${cmd} ${args}"
    ;;
recreate)
    services=$*
    eval "${compose} up -d --no-deps ${services}"
    ;;
run)
    service=$1
    eval "${compose} run --rm ${service}"
    ;;
dbdump)
    eval "${compose} run db_backup /dump.sh"
    ;;
*)
cat << EOF
Usage ./env.sh <environemnt> <command> [args...]

Environments:
    dev  - used for local development.
           Examples:
           ./env.sh dev up platform_devices command_executor
           ./env.sh dev ps
           ./env.sh dev logs
           ./env.sh dev logs command_executor rabbitmq
    prod - used for production deployment.
    test - used for running integration tests on CI.

Commands:
    up, down, recreate
    run, exec
    ps, logs
    start, stop, restart
    dbdump

Config files for specific environment should be placed at ~/.tl/<environment>.env. For example:
    export IMAGE_TAG=v0.6
    export PROXY_PORT=2347
EOF
    exit 1
    ;;
esac
