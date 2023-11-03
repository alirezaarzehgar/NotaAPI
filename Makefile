APP_ENVS := MYSQL_HOST=localhost

all: prod

prod: stop
	go mod vendor
	docker-compose build
	docker-compose up -d

dev: .env
	docker-compose up db -d
	$(APP_ENVS) go run .

.env:
	cp .env.example .evn
	mkdir logs assets

stop:
	docker-compose stop
