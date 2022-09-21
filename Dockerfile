# syntax=docker/dockerfile:1

#
# Build
#
FROM golang:1.19.1-buster AS build
ENV CGO_ENABLED 0
ENV GOOS linux
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" cmd/mool/main.go

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /app

COPY --from=build /app/main /app/main
COPY --from=build ./app/app.env /app/
COPY --from=build /app/grpc_health_probe /app/grpc_health_probe
COPY --from=build /go/bin/dlv /app/dlv

EXPOSE 22210
EXPOSE 22211
EXPOSE 2345

USER nonroot:nonroot

ENTRYPOINT ["/app/dlv"]