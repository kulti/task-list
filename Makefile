build: *.go
	go build -o tl

run-live-reload:
	modd

run: build
	./tl server -p 8090

gen-go:
	docker run --rm -it -v ${PWD}:/local openapitools/openapi-generator-cli generate --package-name=openapi_cli -Dapis,models,supportingFiles=client.go -i /local/api/task.yaml -g go -o /local/internal/router/openapi_cli
	docker run --rm -it -v ${PWD}:/local openapitools/openapi-generator-cli generate --package-name=openapi_cli -DsupportingFiles=configuration.go -i /local/api/task.yaml -g go -o /local/internal/router/openapi_cli

gen-ts:
	docker run --rm -it -v ${PWD}:/local openapitools/openapi-generator-cli generate -i /local/api/task.yaml -g typescript-jquery -o /local/frontend/ts/openapi_cli

gen-css:
	frontend/scss/task_menu.scss frontend/css/main.css

build-js:
	tsc --strict --outDir frontend/js frontend/ts/main.ts && \
	browserify frontend/js/main.js > frontend/js/bundle.js
