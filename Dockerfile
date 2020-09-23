# compile stage
ARG GO_VERSION=1.15.2
FROM golang:${GO_VERSION} as immediate

# build binary
COPY . /repo/cube
WORKDIR /repo/cube

RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -v -a -tags netgo -installsuffix netgo -ldflags "-s -w"

# final docker image building stage
FROM golang:${GO_VERSION} as builder
COPY --from=immediate /repo/cube/cube /app/cube

VOLUME /root/.config/cube/cache

ENTRYPOINT ["/app/cube"]
CMD ["--help"]
