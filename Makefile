gen-go: server/internal/generated/openapicli/api_default.go

server/internal/generated/openapicli/api_default.go: api/task.yaml
	docker run --rm -it -v ${PWD}:/local openapitools/openapi-generator-cli:v4.3.1 generate --package-name=openapicli -Dapis,models,supportingFiles=client.go -i /local/api/task.yaml -g go -o /local/internal/generated/openapicli
	docker run --rm -it -v ${PWD}:/local openapitools/openapi-generator-cli:v4.3.1 generate --package-name=openapicli -DsupportingFiles=configuration.go -i /local/api/task.yaml -g go -o /local/internal/generated/openapicli


gen-ts: frontend/src/openapi_cli/index.ts

frontend/src/openapi_cli/index.ts: api/task.yaml
	docker run --rm -it -v ${PWD}:/local openapitools/openapi-generator-cli:v4.3.1 generate -i /local/api/task.yaml -g typescript-jquery -o /local/frontend/src/openapi_cli

build-js: gen-ts frontend/dist/bundle.js

frontend/dist/bundle.js: frontend/src/index.ts
	cd frontend && \
	npx webpack

build-docker-tl-proxy:
	DOCKER_BUILDKIT=1 docker build -f build/package/proxy.Dockerfile -t tl-proxy ./proxy

build-docker-tl-server:
	DOCKER_BUILDKIT=1 docker build -f build/package/tl-server.Dockerfile -t tl-server ./server

build-docker-tl-front: build-js
	DOCKER_BUILDKIT=1 docker build -f build/package/tl-front.Dockerfile -t tl-front ./frontend

build-docker-tl-migrate:
	DOCKER_BUILDKIT=1 docker build -f build/package/tl-migrate.Dockerfile -t tl-migrate ./db

build-docker-tl-integration-tests: build-docker-tl-server build-docker-tl-migrate
	DOCKER_BUILDKIT=1 docker build -f build/package/tl-integration-tests.Dockerfile -t tl-integration-tests ./server

run-tl-prod: build-docker-tl-proxy build-docker-tl-front build-docker-server build-docker-tl-migrate
	cd deployments && \
	docker-compose -p prod -f docker-compose.yaml -f docker-compose.prod.yaml up

stop-tl-prod:
	cd deployments && \
	docker-compose -p prod -f docker-compose.yaml -f docker-compose.prod.yaml down

run-tl-integration-tests: build-docker-tl-integration-tests
	cd deployments && \
	docker-compose -p integration-tests -f docker-compose.yaml -f docker-compose.tests.yaml run db_migrations up && \
	docker-compose -p integration-tests -f docker-compose.yaml -f docker-compose.tests.yaml run tl-integration-tests; \
	docker-compose -p integration-tests -f docker-compose.yaml -f docker-compose.tests.yaml down

build-docker-modd:
	DOCKER_BUILDKIT=1 docker build -f build/package/modd.Dockerfile -t tl-live-reload .

run-tl-dev: build-docker-tl-proxy build-docker-tl-front build-docker-modd build-docker-tl-migrate
	export SRC=${PWD}; \
	cd deployments && \
	docker-compose -p dev -f docker-compose.yaml -f docker-compose.dev.yaml run db_migrations up && \
	docker-compose -p dev -f docker-compose.yaml -f docker-compose.dev.yaml up

stop-tl-dev:
	export SRC=${PWD}; \
	cd deployments && \
	docker-compose -p dev -f docker-compose.yaml -f docker-compose.dev.yaml down

db-dump:
	cd deployments && \
	source database.env && \
	docker-compose -p prod -f docker-compose.yaml -f docker-compose.prod.yaml exec db pg_dump --username=$$POSTGRES_USER --dbname=$$POSTGRES_DB --data-only > db.dump && \
	sed -i '' -e 's/COPY public.task_lists /DELETE FROM public.task_lists;\'$$'\nCOPY public.task_lists /' \
		-e '/COPY public.schema_migrations /{N;N;d;}' \
	db.dump

db-restore:
	cd deployments && \
	source database.env && \
	docker-compose -p prod -f docker-compose.yaml -f docker-compose.prod.yaml run db_migrations up && \
	db_container=$(docker-compose -p prod -f docker-compose.yaml -f docker-compose.prod.yaml ps -q db) && \
	docker cp db.dump ${db_container}:/tmp/ && \
	docker-compose -p prod -f docker-compose.yaml -f docker-compose.prod.yaml exec -T db psql --username=$$POSTGRES_USER --dbname=$$POSTGRES_DB -f /tmp/db.dump

go-coverage:
	cd server && \
	./scripts/go_test.sh && \
	go tool cover -html=coverage.txt
