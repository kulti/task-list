#!/bin/bash

set -e

cd deployments
source tests.env

trap "docker-compose -p integration-tests -f docker-compose.yaml -f docker-compose.tests.yaml down" EXIT

docker-compose -p integration-tests -f docker-compose.yaml -f docker-compose.tests.yaml run db_migrations up
docker-compose -p integration-tests -f docker-compose.yaml -f docker-compose.tests.yaml run tl-integration-tests
