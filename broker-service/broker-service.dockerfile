FROM golang:1.18-alpine as builder

# run CMD create folder app
RUN mkdir /app

# Copy all things in current folder to /app
COPY . /app

# Set working directory
WORKDIR /app

# Do not use C lib
RUN CGO_ENABLE=0 go build -o brokerApp ./cmd/api

# Set executable
RUN chmod +x /app/brokerApp

# Build tiny docker images
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/brokerApp /app

# Executive with escapse string
CMD ["/app/brokerApp"]
