APP_ENVS := MYSQL_HOST=localhost

all: prod

prod: stop
	go mod vendor
	docker-compose build
	docker-compose up -d

dev:
	docker-compose up db -d
	$(APP_ENVS) go run .

stop:
	docker-compose stop
