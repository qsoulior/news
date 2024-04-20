FROM golang:1.22-alpine3.19 AS dependencies
WORKDIR /dependencies
COPY ./aggregator ./aggregator
COPY ./parser ./parser
COPY ./lenta-parser/go.mod ./lenta-parser/go.sum ./lenta-parser/
WORKDIR /dependencies/lenta-parser
RUN go mod download && go mod verify

FROM dependencies AS build
WORKDIR /build
COPY ./lenta-parser ./lenta-parser
COPY --from=dependencies /dependencies/aggregator ./aggregator
COPY --from=dependencies /dependencies/parser ./parser
WORKDIR /build/lenta-parser
RUN CGO_ENABLED=0 GOOS=linux go build -v -o ../main ./cmd

FROM alpine:3.19
WORKDIR /app
COPY --from=build /build/main ./
ENTRYPOINT ["./main"]