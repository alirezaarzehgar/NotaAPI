APP_ENVS := MYSQL_HOST=localhost

all: prod

prod: stop .env
	go mod vendor
	docker-compose build
	docker-compose up -d

dev: .env
	docker-compose up db
	$(APP_ENVS) go run .

.env:
	cp .env.example .env
	mkdir logs assets

stop:
	docker-compose stop
