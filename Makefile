gen-go: server/internal/generated/openapicli/api_default.go

server/internal/generated/openapicli/api_default.go: api/task.yaml
	docker run --rm -it -v ${PWD}:/local openapitools/openapi-generator-cli:v4.3.1 generate --package-name=openapicli -Dapis,models,supportingFiles=client.go -i /local/api/task.yaml -g go -o /local/internal/generated/openapicli
	docker run --rm -it -v ${PWD}:/local openapitools/openapi-generator-cli:v4.3.1 generate --package-name=openapicli -DsupportingFiles=configuration.go -i /local/api/task.yaml -g go -o /local/internal/generated/openapicli


gen-ts: front/src/openapi_cli/index.ts

front/src/openapi_cli/index.ts: api/task.yaml
	docker run --rm -it -v ${PWD}:/local openapitools/openapi-generator-cli:v4.3.1 generate -i /local/api/task.yaml -g typescript-jquery -o /local/front/src/openapi_cli

build-js: gen-ts front/dist/bundle.js

front/dist/bundle.js: front/src/index.ts
	cd front && \
	npx webpack

SERVICES=proxy server front migrate live-reload

front: build-js

$(addprefix build-docker-tl-, $(SERVICES)): build-docker-tl-%: %
	DOCKER_BUILDKIT=1 docker build -t tl-$< ./$<

build-docker-tl-integration-tests:
	DOCKER_BUILDKIT=1 docker build -f server/tl-integration-tests.Dockerfile -t tl-integration-tests ./server

run-tl-integration-tests: build-docker-tl-integration-tests
	cd deployments && \
	docker-compose -p integration-tests -f docker-compose.yaml -f docker-compose.tests.yaml run db_migrations up && \
	docker-compose -p integration-tests -f docker-compose.yaml -f docker-compose.tests.yaml run tl-integration-tests; \
	docker-compose -p integration-tests -f docker-compose.yaml -f docker-compose.tests.yaml down

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

include environments.mk
