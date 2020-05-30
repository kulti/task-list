gen-go:
	docker run --rm -it -v ${PWD}:/local openapitools/openapi-generator-cli generate --package-name=openapicli -Dapis,models,supportingFiles=client.go -i /local/api/task.yaml -g go -o /local/internal/generated/openapicli
	docker run --rm -it -v ${PWD}:/local openapitools/openapi-generator-cli generate --package-name=openapicli -DsupportingFiles=configuration.go -i /local/api/task.yaml -g go -o /local/internal/generated/openapicli

gen-ts:
	docker run --rm -it -v ${PWD}:/local openapitools/openapi-generator-cli generate -i /local/api/task.yaml -g typescript-jquery -o /local/frontend/ts/openapi_cli

gen-css:
	frontend/scss/task_menu.scss frontend/css/main.css

build-js:
	tsc --strict --outDir frontend/js frontend/ts/main.ts && \
	browserify frontend/js/main.js > frontend/js/bundle.js

build-docker-tl-server:
	DOCKER_BUILDKIT=1 docker build -f build/package/tl-server.Dockerfile -t tl-server .

build-docker-tl-migrate:
	DOCKER_BUILDKIT=1 docker build -f build/package/tl-migrate.Dockerfile -t tl-migrate ./db

build-docker-tl-integration-tests: build-docker-tl-server build-docker-tl-migrate
	DOCKER_BUILDKIT=1 docker build -f build/package/tl-integration-tests.Dockerfile -t tl-integration-tests .

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
	docker-compose -p dev -f docker-compose.yaml -f docker-compose.dev.yaml run --service-ports tl_live_reload

stop-tl-dev:
	cd deployments && \
	docker-compose -p dev -f docker-compose.yaml -f docker-compose.dev.yaml down
