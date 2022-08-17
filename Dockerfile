# syntax=docker/dockerfile:1

##
## Build
##
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
# EXPOSE 32333

# USER nonroot:nonroot

# ENTRYPOINT ["/app/main"]


FROM golang:1.17-buster
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go get -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" cmd/mool/main.go
COPY /go/bin/dlv /app/dlv
EXPOSE 22211
EXPOSE 32333
ENTRYPOINT [ "/app/dlv" , "--listen=:32333", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/app/main"]]