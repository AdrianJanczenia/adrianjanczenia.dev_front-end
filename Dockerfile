FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/frontend ./main.go

# ---

FROM gcr.io/distroless/static-debian11

WORKDIR /app

COPY --from=builder /app/config ./config
COPY --from=builder /app/internal/web ./internal/web
COPY --from=builder /app/frontend .

USER nonroot:nonroot

CMD ["/app/frontend"]