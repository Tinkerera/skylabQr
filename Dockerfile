# Build aşaması
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o url-shortener

# Çalıştırma aşaması
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/url-shortener .
EXPOSE 8080
CMD ["sh", "-c", "sleep 10 && ./url-shortener"]
