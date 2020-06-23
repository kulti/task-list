#!/bin/bash

set -e

env=$1
action=$2

function up() {
    env=$1
    docker-compose -p ${env} -f docker-compose.yaml -f docker-compose.${env}.yaml run db_migrations up
	docker-compose -p ${env} -f docker-compose.yaml -f docker-compose.${env}.yaml up -d
}

function cmd() {
    env=$1
    action=$2
	docker-compose -p ${env} -f docker-compose.yaml -f docker-compose.${env}.yaml ${action}
}

cd deployments
[ -f ${env}.env ] && source ${env}.env
[ -f ~/.tl/${env}.env ] && source ~/.tl/${env}.env

case $action in
up)
    up ${env}
    ;;
ps|down)
    cmd ${env} ${action}
    ;;
restart)
    cmd ${env} down
    up ${env}
    ;;
logs)
    cmd ${env} "logs -f"
    ;;
*)
cat << EOF
Usage ./scripts/deploy.sh environemnt command

Commands:
    up
    down
    ps
    logs

Environments:
    dev
    prod

Config files for specific environment should be placed at ~/.tl/<environment>.env. For example:
    export IMAGE_TAG=v0.6
    export PROXY_PORT=2347
EOF
    ;;
esac
