.PHONY: ratelimiter.local
ratelimiter.local:
	go run ./cmd/app

.PHONY: docker.build
docker.build:
	docker build -f ./deployments/build/Dockerfile -t ratelimiter .

.PHONY: docker.deploy
docker.deploy:
	docker build -f ./deployments/build/Dockerfile -t ratelimiter . && \
	docker-compose -p ratelimiter -f ./deployments/environment/docker-compose.deploy.yml up -d

.PHONY: docker.deploy.down
docker.deploy.down:
	docker-compose -p ratelimiter -f ./deployments/environment/docker-compose.deploy.yml down

.PHONY: ratelimiter.dev.env
ratelimiter.dev.env:
	docker-compose -p ratelimiter -f ./deployments/environment/docker-compose.dev.yml up -d

.PHONY: ratelimiter.dev.env.down
ratelimiter.dev.env.down:
	docker-compose -p ratelimiter -f ./deployments/environment/docker-compose.dev.yml down

.PHONY: uint.testing
uint.testing:
	go test ./...

.PHONY: integration.testing
integration.testing:
	docker-compose -p ratelimiter -f ./deployments/environment/docker-compose.dev.yml up -d && \
	go clean -testcache && go test ./cmd/app && \
	docker-compose -p ratelimiter -f ./deployments/environment/docker-compose.dev.yml down