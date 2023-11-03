# NotaAPI

## Run appliaction using docker-compose
First of all write your `.env` file.
```
cp .env.example .env
```

For development environemnt
```
make dev
```
or 
```
MYSQL_HOST=localhost go run .
```

You can run your MySQL server and pass it's hostname to `MYSQL_HOST` environment variable.
Even you can run a MySQL server using following command:
```
docker-compose up db -d
```
Run on production:
```
make prod
```
### Info
Default web server port number: 8000
You can change it with `RUNNING_PORT` environment variable on `.env` or your shell.
