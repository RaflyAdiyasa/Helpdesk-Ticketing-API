FROM golang:1.25.5-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /bin/helpdesk-api ./cmd/api

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app

RUN apk add --no-cache curl ca-certificates

WORKDIR /app

COPY --from=builder /bin/helpdesk-api /app/helpdesk-api

EXPOSE 8080

USER app

ENTRYPOINT ["/app/helpdesk-api"]
