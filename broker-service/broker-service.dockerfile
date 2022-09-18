FROM alpine:latest

RUN mkdir /app

COPY brokerApp /app

# Executive with escapse string
CMD ["/app/brokerApp"]
