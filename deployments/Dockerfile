# syntax=docker/dockerfile:1

FROM node:20-alpine AS frontend-builder
RUN corepack enable && corepack prepare pnpm@10.15.0 --activate
WORKDIR /frontend
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY frontend/ ./
RUN pnpm build

FROM golang:1.25-alpine AS backend-builder
WORKDIR /backend
COPY backend/go.mod backend/go.sum ./
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o shiyi ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=backend-builder /backend/shiyi /app/shiyi
COPY --from=frontend-builder /frontend/dist /app/public
COPY backend/migrations /app/migrations
ENV APP_PORT=8080 \
    TZ=Asia/Shanghai \
    MIGRATION_DIR=/app/migrations \
    PUBLIC_DIR=/app/public \
    UPLOAD_DIR=/app/uploads
EXPOSE 8080
CMD ["/app/shiyi"]
