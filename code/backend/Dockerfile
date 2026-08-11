# syntax=docker/dockerfile:1

# The build context is this directory — see `build.context` in
# docker-compose.yml.

FROM golang:1.25-alpine AS builder

WORKDIR /src

# Dependencies first so a code-only change does not re-download the module
# cache.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ./... rather than a named package path, which is the same command the
# workspace runtime uses to build a project. It requires the module to hold
# exactly one main package — a constraint that already applied before this file
# existed, and one worth keeping: two entry points and neither runtime knows
# which one is the server.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /bin/server ./...

# ── Runtime ──────────────────────────────────────────────────────────────────

FROM alpine:3.20

# ca-certificates for outbound TLS, tzdata so timestamps are not all UTC by
# accident, wget for the healthcheck in docker-compose.yml.
RUN apk add --no-cache ca-certificates tzdata wget

WORKDIR /app

COPY --from=builder /bin/server /app/server

ENV PORT=8080

USER nobody

EXPOSE 8080

CMD ["/app/server"]
