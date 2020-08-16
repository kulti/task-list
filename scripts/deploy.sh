#!/bin/bash

set -e

env=$1
action=$2
extra_args=$3

echo "~${extra_args}~"

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

function gen_migration_name() {
    env=$1
    name=$2
    echo "~docker-compose -p ${env} -f docker-compose.yaml -f docker-compose.${env}.yaml run db_migrations create ${name}~"
    docker-compose -p ${env} -f docker-compose.yaml -f docker-compose.${env}.yaml run db_migrations create -ext sql ${name}
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
dbdump)
    cmd ${env} "run db_backup /dump.sh"
    ;;
gen-migration-name)
    gen_migration_name ${env} ${extra_args}
    ;;
*)
cat << EOF
Usage ./scripts/deploy.sh environemnt command

Commands:
    up
    down
    ps
    logs
    dbdump

Environments:
    dev
    prod

Config files for specific environment should be placed at ~/.tl/<environment>.env. For example:
    export IMAGE_TAG=v0.6
    export PROXY_PORT=2347
EOF
    ;;
esac
