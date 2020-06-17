FROM golang:1.14 AS builder

WORKDIR /app

RUN go mod init build.com/build && \
    CGO_ENABLED=0 go get -v -ldflags='-w -s -extldflags "-static"' github.com/cortesi/modd/cmd/modd

FROM golang:1.14-alpine

COPY --from=builder /go/bin/modd /modd

WORKDIR /app
ENTRYPOINT [ "/modd" ]
