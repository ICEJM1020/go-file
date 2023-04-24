ARG VERSION V2.2

# Build stage
FROM golang AS builder
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -ldflags "-s -w -X 'go-file/common.Version=$(VERSION)' -extldflags '-static'" -o go-file

# Get Aria2 Pro
# FROM p3terx/s6-alpine
FROM alpine

RUN apk update \
    && apk upgrade \
    && apk add --no-cache  jq findutils ca-certificates tzdata \
    && update-ca-certificates 2>/dev/null || true \
    && rm -rf /var/cache/apk/* /tmp/*

ENV PORT=3000 \
    UPDATE_TRACKERS=true \
    CUSTOM_TRACKER_URL= \
    LISTEN_PORT=6888 \
    RPC_PORT=6800 \
    RPC_SECRET=TIMBER \
    DISK_CACHE=64M\
    UMASK_SET= \
    SPECIAL_MODE=

COPY --from=builder /build/go-file /
COPY --from=builder /build/bin/aria2c /

RUN chmod u+x /aria2c

WORKDIR /data
EXPOSE 3000 \
    6800 \
    6888 \
    6888/udp
VOLUME [ "/data" ]
ENTRYPOINT ["/go-file"]
