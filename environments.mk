ENVIRONMENTS=prod dev

up-tl-prod-deps: build-docker-tl-proxy build-docker-tl-front build-docker-tl-server build-docker-tl-migrate
up-tl-dev-deps: build-docker-tl-proxy build-docker-tl-front build-docker-tl-live-reload build-docker-tl-migrate

$(addprefix up-tl-, $(ENVIRONMENTS)): up-tl-%: % up-tl-%-deps
	./deployments/env.sh $< up

$(addprefix down-tl-, $(ENVIRONMENTS)): down-tl-%: %
	./deployments/env.sh $< down

$(addprefix ps-tl-, $(ENVIRONMENTS)): ps-tl-%: %
	./deployments/env.sh $< ps

$(addprefix logs-tl-, $(ENVIRONMENTS)): logs-tl-%: %
	./deployments/env.sh $< logs

.PHONY: $(ENVIRONMENTS)
