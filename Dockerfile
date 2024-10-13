FROM golang:1.23-alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o /bin/ledgernet ./cmd

FROM gcr.io/distroless/static-debian11

COPY --from=builder /bin/ledgernet /bin/ledgernet

USER nonroot:nonroot

EXPOSE 50051

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD [ "grpc_health_probe", "-addr=localhost:50051" ] || exit 1

CMD ["/bin/ledgernet"]
