# syntax=docker/dockerfile:1

#
# Build
#
FROM golang:1.17-buster AS build
ENV CGO_ENABLED 0
ENV GOOS linux
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/mool/main.go

# Adding the grpc_health_probe
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.11 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /app

COPY --from=build /app/main /app/main
COPY --from=build ./app/app.env /app/
COPY --from=build /bin/grpc_health_probe /app/grpc_health_probe

EXPOSE 22210

USER nonroot:nonroot

ENTRYPOINT ["/app/main"]