# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17-buster AS build
ENV CGO_ENABLED 0
ENV GOOS linux
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# COPY *.go ./
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build /cmd/mool/main.go -o /mool

##
## Deploy
##
FROM gcr.io/distroless/base-debian10
# FROM gcr.io/distroless/static-debian10:nonroot

WORKDIR /app

COPY --from=build /mool /app/mool
COPY --from=build ./app/app.env /app/

EXPOSE 22210

USER nonroot:nonroot

ENTRYPOINT ["/app/mool"]