FROM golang as builder

ADD . /src

WORKDIR /src

ENV GO111MODULE on
RUN go get -d -v ./...
RUN go install -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o urlshortener .

FROM scratch

ENV PORT 8080

COPY --from=builder /src /app/

WORKDIR /app

CMD ["./urlshortener"]