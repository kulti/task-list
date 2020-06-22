ENVIRONMENTS=prod dev

up-tl-prod-deps: build-docker-tl-proxy build-docker-tl-front build-docker-tl-server build-docker-tl-migrate
up-tl-dev-deps: build-docker-tl-proxy build-docker-tl-front build-docker-tl-live-reload build-docker-tl-migrate

$(addprefix up-tl-, $(ENVIRONMENTS)): up-tl-%: % up-tl-%-deps
	./scripts/deploy.sh $< up

$(addprefix down-tl-, $(ENVIRONMENTS)): down-tl-%: %
	./scripts/deploy.sh $< down

$(addprefix ps-tl-, $(ENVIRONMENTS)): ps-tl-%: %
	./scripts/deploy.sh $< ps

$(addprefix logs-tl-, $(ENVIRONMENTS)): logs-tl-%: %
	./scripts/deploy.sh $< logs

.PHONY: $(ENVIRONMENTS)
