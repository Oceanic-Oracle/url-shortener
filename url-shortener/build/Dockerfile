FROM golang:alpine AS build

RUN mkdir /app

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY internal ./internal
COPY cmd ./cmd

RUN go build -o /app/bin ./cmd/shortener/main.go

FROM alpine:3.21

RUN apk --no-cache add curl

COPY --from=build /app/bin /app/bin

EXPOSE 8080 8081

ENTRYPOINT ["/app/bin"]