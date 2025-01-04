FROM golang:1.23 AS builder

WORKDIR /usr/src/app

RUN GRPC_HEALTH_PROBE_VERSION=v0.4.13 && \
  wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
  chmod +x /bin/grpc_health_probe

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading
# them in subsequent builds if they change.
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -o /go/bin/app ./cmd

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder --chown=nonroot:nonroot /go/bin/app /
COPY --from=builder --chown=nonroot:nonroot /bin/grpc_health_probe /

EXPOSE 8080/tcp

CMD ["/app", "run"]

