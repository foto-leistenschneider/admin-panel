FROM golang:1-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /app/main .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/main /app/main
CMD ["/app/main"]
