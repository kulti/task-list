gen-go:
	docker run --rm -it -v ${PWD}:/local openapitools/openapi-generator-cli generate --package-name=openapicli -Dapis,models,supportingFiles=client.go -i /local/api/task.yaml -g go -o /local/internal/generated/openapicli
	docker run --rm -it -v ${PWD}:/local openapitools/openapi-generator-cli generate --package-name=openapicli -DsupportingFiles=configuration.go -i /local/api/task.yaml -g go -o /local/internal/generated/openapicli

gen-ts:
	docker run --rm -it -v ${PWD}:/local openapitools/openapi-generator-cli generate -i /local/api/task.yaml -g typescript-jquery -o /local/frontend/src/openapi_cli

build-js:
	cd frontend && \
	npx webpack

build-docker-tl-proxy:
	DOCKER_BUILDKIT=1 docker build -f build/package/proxy.Dockerfile -t tl-proxy ./proxy

build-docker-tl-server:
	DOCKER_BUILDKIT=1 docker build -f build/package/tl-server.Dockerfile -t tl-server .

build-docker-tl-front:
	DOCKER_BUILDKIT=1 docker build -f build/package/tl-front.Dockerfile -t tl-front ./frontend

build-docker-tl-migrate:
	DOCKER_BUILDKIT=1 docker build -f build/package/tl-migrate.Dockerfile -t tl-migrate ./db

build-docker-tl-integration-tests: build-docker-tl-server build-docker-tl-migrate
	DOCKER_BUILDKIT=1 docker build -f build/package/tl-integration-tests.Dockerfile -t tl-integration-tests .

run-tl-prod:
	cd deployments && \
	docker-compose -p task-list -f docker-compose.yaml -f docker-compose.prod.yaml up

run-tl-integration-tests: build-docker-tl-integration-tests
	cd deployments && \
	docker-compose -p integration-tests -f docker-compose.yaml -f docker-compose.tests.yaml run db_migrations up && \
	docker-compose -p integration-tests -f docker-compose.yaml -f docker-compose.tests.yaml run tl-integration-tests; \
	docker-compose -p integration-tests -f docker-compose.yaml -f docker-compose.tests.yaml down

build-docker-modd:
	DOCKER_BUILDKIT=1 docker build -f build/package/modd.Dockerfile -t tl-live-reload .

run-tl-dev:
	export SRC=${PWD}; \
	cd deployments && \
	docker-compose -p dev -f docker-compose.yaml -f docker-compose.dev.yaml run db_migrations up && \
	docker-compose -p dev -f docker-compose.yaml -f docker-compose.dev.yaml up

stop-tl-dev:
	export SRC=${PWD}; \
	cd deployments && \
	docker-compose -p dev -f docker-compose.yaml -f docker-compose.dev.yaml down
