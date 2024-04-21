FROM golang:1.22-alpine3.19 AS dependencies
WORKDIR /dependencies

COPY ./aggregator/entity ./aggregator/entity
COPY ./aggregator/pkg/rabbitmq ./aggregator/pkg/rabbitmq
COPY ./aggregator/go.mod ./aggregator/go.sum ./aggregator/
COPY ./parser ./parser
COPY ./ria-parser/go.mod ./ria-parser/go.sum ./ria-parser/

WORKDIR /dependencies/ria-parser
RUN go mod download && go mod verify

FROM golang:1.22-alpine3.19 AS build
WORKDIR /build
COPY ./ria-parser ./ria-parser
COPY --from=dependencies /dependencies/aggregator ./aggregator
COPY --from=dependencies /dependencies/parser ./parser
WORKDIR /build/ria-parser
RUN CGO_ENABLED=0 GOOS=linux go build -v -o ../main ./cmd

FROM ghcr.io/go-rod/rod:v0.115.0
WORKDIR /app
COPY --from=build /build/main ./
ENTRYPOINT ["./main"]