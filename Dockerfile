FROM golang:1.17 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/cli/main.go

# Multistage because we can just run the binary we need without using the golang image which is much larger and contains lots
# of tools that are useless during runtime
FROM alpine:latest AS production 
RUN mkdir /go-github-cli
COPY --from=builder /app /go-github-cli
WORKDIR /go-github-cli
ENTRYPOINT ["./app"]