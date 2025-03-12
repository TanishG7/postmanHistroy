FROM golang:latest AS builder

WORKDIR /www/var/GOSERVER

COPY ./GOSERVER .

RUN ["go", "mod", "tidy"]

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o build .

FROM debian

RUN apt-get update && apt-get install -y --no-install-recommends util-linux && rm -rf /var/lib/apt/lists/*

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /www

COPY --from=builder /www/var/GOSERVER .

RUN chmod 777 ./build

CMD ["/bin/sh", "-c", "./build > ./logfile.log 2>&1"]