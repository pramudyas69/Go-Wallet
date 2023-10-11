# Stage 1: Build the Go application
FROM golang:1.16 AS builder

WORKDIR /app
COPY . .
RUN go build -o myapp

FROM alpine:3.14

COPY --from=builder /app/myapp /usr/local/bin/myapp
WORKDIR /app
EXPOSE 8080
CMD ["myapp"]



