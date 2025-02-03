FROM golang:latest AS builder

WORKDIR /www/var/GOSERVER

COPY ./GOSERVER .

RUN ["go", "mod", "tidy"]

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o build .

FROM alpine:latest

WORKDIR /www

COPY --from=builder /www/var/GOSERVER . 

RUN ["chmod", "777", "./build"]

CMD ["./build", ">", "./logfile.log", "2>&1"]
