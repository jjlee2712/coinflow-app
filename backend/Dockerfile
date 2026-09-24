FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server .

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]