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

## Usage
Encode any URL by running:
```
curl --location --request POST 'localhost:8080/encode' \
--header 'access_token: test' \
--header 'Content-Type: application/json' \
--data-raw '{
    "url": "https://www.google.com/"
}'
```

If successful, the service will return a response similar to the following:
```
{
    "success": true,
    "data": "localhost:8080/tE5yk"
}
```
If you then access the URL in "data", you will be redirected to the URL from your original request.