FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o roleleader cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app/

COPY --from=build /app/roleleader .
COPY /config/prodConfig.yaml /app/config/config.yaml
COPY /storage/migrations /app/storage/migrations

CMD ["./roleleader"]