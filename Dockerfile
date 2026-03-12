# --- VPSBENCH (vpsbench) ---
# Multi-stage Dockerfile для CLI-инструмента бенчмаркинга

# ============================================================
# Stage 1: Builder — сборка бинарника
# ============================================================
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build

# Кэширование зависимостей
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Сборка
COPY . .
ARG VERSION=dev
ARG COMMIT=unknown

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build \
    -ldflags="-s -w" \
    -o /vpsbench ./cmd/vpsbench/

# ============================================================
# Stage 2: Development — разработка с hot reload
# ============================================================
FROM golang:1.24-alpine AS development

RUN apk add --no-cache git

# Установка air для hot reload
RUN go install github.com/air-verse/air@latest

WORKDIR /app
COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

CMD ["air"]

# ============================================================
# Stage 3: Production — минимальный образ
# ============================================================
FROM alpine:3.21 AS production

LABEL org.opencontainers.image.title="vpsbench" \
      org.opencontainers.image.description="Comprehensive VPS benchmarking tool" \
      org.opencontainers.image.source="https://github.com/user/vpsbench" \
      org.opencontainers.image.vendor="VPSBench"

# Сетевые утилиты нужны для бенчмарка сети
RUN apk add --no-cache ca-certificates curl

# Non-root пользователь
RUN addgroup -S bench && adduser -S bench -G bench
USER bench

COPY --from=builder /vpsbench /usr/local/bin/vpsbench

ENTRYPOINT ["vpsbench"]
