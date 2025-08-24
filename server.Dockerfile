# syntax=docker/dockerfile:1.7-labs

FROM golang:alpine AS builder
WORKDIR /app
COPY  --exclude=client/ . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o app server/server.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 7070
CMD ["./app"]
