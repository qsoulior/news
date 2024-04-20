FROM golang:1.22-alpine3.19 AS dependencies
WORKDIR /dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

FROM dependencies AS build
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o ./main ./cmd

FROM alpine:3.19
WORKDIR /app
COPY --from=build /build/main ./
ENTRYPOINT ["./main"]