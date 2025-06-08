FROM node:22-alpine AS styles
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
WORKDIR /app
COPY package.json pnpm-lock.yaml .
COPY ./internal/ ./internal/
RUN pnpm install --frozen-lockfile && \
    pnpm run styles

FROM golang:1-alpine AS binary
WORKDIR /app
COPY . .
COPY --from=styles /app/internal/server/view/styles.min.css /app/internal/server/view/styles.min.css
RUN go install github.com/a-h/templ/cmd/templ@latest && \
    go mod download && \
    templ generate && \
    go build -o /app/main .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=binary /app/main /app/main
CMD ["/app/main"]
