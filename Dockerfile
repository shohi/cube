# compile stage
ARG GO_VERSION=1.12.7
FROM golang:${GO_VERSION} as immediate

# download the source
WORKDIR /go/src/github.com/shohi/cube
RUN git clone https://github.com/shohi/cube.git .

# binary build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -tags netgo -installsuffix netgo -ldflags "-s -w"

# final docker image building stage
FROM golang:${GO_VERSION} as builder
COPY --from=immediate /go/src/github.com/shohi/cube/cube /app/cube 

VOLUME /root/.config/cube/cache

ENTRYPOINT ["/app/cube"]
CMD ["--help"]
