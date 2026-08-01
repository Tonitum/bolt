FROM golang:1.25.2-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -o /bolt ./cmd/core-service

FROM alpine:3.22

WORKDIR /data

COPY --from=build /bolt /usr/local/bin/bolt
COPY --from=build /src/static /static

VOLUME ["/data"]
EXPOSE 80

ENTRYPOINT ["/usr/local/bin/bolt"]
