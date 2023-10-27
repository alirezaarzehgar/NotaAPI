APP_ENVS := DB_HOSTNAME=localhost

all: prod

prod: stop
	go mod vendor
	docker-compose build
	docker-compose up

dev:
	docker-compose up db -d
	$(APP_ENVS) go run .

stop:
	docker-compose stop
