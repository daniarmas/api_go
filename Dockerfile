# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17-buster AS build
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# COPY *.go ./
COPY . .

RUN go build -o /api_go

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /app

COPY --from=build /api_go /app/api_go
COPY --from=build ./app/app.env /app/

EXPOSE 22210

USER nonroot:nonroot

ENTRYPOINT ["/app/api_go"]