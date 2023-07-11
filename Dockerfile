FROM golang:1.19 AS build

LABEL maintainer="Saman Hoseini <saman2000hoseini@gmail.com>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o webhook-proxy ./cmd/webhook-proxy

#Second stage of build
FROM alpine:3.18

RUN apk --no-cache add ca-certificates

WORKDIR /app/

COPY --from=build /app/webhook-proxy .

EXPOSE 8080

CMD ["/app/webhook-proxy", "api"]
