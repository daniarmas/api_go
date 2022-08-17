# syntax=docker/dockerfile:1

# ##
# ## Build
# ##
# FROM golang:1.17-buster AS build
# ENV CGO_ENABLED 0
# ENV GOOS linux
# WORKDIR /app

# COPY go.mod ./
# COPY go.sum ./
# RUN go mod download

# COPY . .

# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/mool/main.go

# ##
# ## Deploy
# ##
# FROM gcr.io/distroless/base-debian10

# WORKDIR /app

# COPY --from=build /app/main /app/main
# COPY --from=build ./app/app.env /app/

# EXPOSE 22210

# USER nonroot:nonroot

# ENTRYPOINT ["/app/main"]

FROM golang:1.17-buster AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go get -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/mool/main.go

FROM gcr.io/distroless/base-debian10
WORKDIR /app
COPY . .
COPY --from=build /app/main /app/main
COPY --from=build ./app/app.env /app/
COPY --from=build /go/bin/dlv /app/
ENTRYPOINT [ "/dlv" , "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/app"]]