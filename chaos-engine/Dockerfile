# Build Stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o chaos-engine main.go

# Run Stage (Tiny image)
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/chaos-engine .

# Command to run the app
CMD ["./chaos-engine"]
