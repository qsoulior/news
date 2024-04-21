FROM golang:1.22-alpine3.19 AS dependencies
WORKDIR /dependencies 

COPY ./aggregator/entity ./aggregator/entity
COPY ./aggregator/pkg/rabbitmq ./aggregator/pkg/rabbitmq
COPY ./aggregator/go.mod ./aggregator/go.sum ./aggregator/
COPY ./parser ./parser
COPY ./iz-parser/go.mod ./iz-parser/go.sum ./iz-parser/

WORKDIR /dependencies/iz-parser
RUN go mod download && go mod verify

FROM golang:1.22-alpine3.19 AS build
WORKDIR /build
COPY ./iz-parser ./iz-parser
COPY --from=dependencies /dependencies/aggregator ./aggregator
COPY --from=dependencies /dependencies/parser ./parser
WORKDIR /build/iz-parser
RUN CGO_ENABLED=0 GOOS=linux go build -v -o ../main ./cmd

FROM alpine:3.19
WORKDIR /app
COPY --from=build /build/main ./
ENTRYPOINT ["./main"]