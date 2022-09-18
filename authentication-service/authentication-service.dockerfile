FROM alpine:latest

RUN mkdir /app

COPY authApp /app

# Executive with escapse string
CMD ["/app/authApp"]
