ENVIRONMENTS=prod dev

up-tl-prod-deps: build-docker-tl-proxy build-docker-tl-front build-docker-tl-server build-docker-tl-migrate
up-tl-dev-deps: build-docker-tl-proxy build-docker-tl-front build-docker-tl-live-reload build-docker-tl-migrate

$(addprefix up-tl-, $(ENVIRONMENTS)): up-tl-%: % up-tl-%-deps
	cd deployments && \
	source $<.env && \
	docker-compose -p $< -f docker-compose.yaml -f docker-compose.$<.yaml run db_migrations up && \
	docker-compose -p $< -f docker-compose.yaml -f docker-compose.$<.yaml up -d

$(addprefix down-tl-, $(ENVIRONMENTS)): down-tl-%: %
	cd deployments && \
	source $<.env && \
	docker-compose -p $< -f docker-compose.yaml -f docker-compose.$<.yaml down

$(addprefix ps-tl-, $(ENVIRONMENTS)): ps-tl-%: %
	cd deployments && \
	source $<.env && \
	docker-compose -p $< -f docker-compose.yaml -f docker-compose.$<.yaml ps

$(addprefix logs-tl-, $(ENVIRONMENTS)): logs-tl-%: %
	cd deployments && \
	source $<.env && \
	docker-compose -p $< -f docker-compose.yaml -f docker-compose.$<.yaml logs -f

.PHONY: $(ENVIRONMENTS)
