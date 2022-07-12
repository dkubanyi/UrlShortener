# URL shortener and redirect service

## Installation
### Using Docker

The simplest way to start the service is docker:

Clone the repo and run docker-compose:

```shell
docker-compose up -d
```

### Locally
You can also use the bot locally by running:

```shell
go build
go run main.go
```
However, it is presumed that you already have a running database, and you put its connection string into the `.env` file.
