### Build stage (make image smaller with AS builder)
FROM golang:1.25.3-alpine3.22 AS builder
WORKDIR /app
# Copy all of current files to current working directory in container
COPY . .
# Build executable file for this package
RUN go build -o main main.go

### Run Stage (makes image smaller with just executable file)
# Goes from ~500MB to ~20MB
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/main .
COPY db/migration ./db/migration

# Container listens on this port (doesn't actually publish the port, just documents it)
EXPOSE 8080
