FROM golang:1.14-alpine AS builder

RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -mod=vendor -o /go/bin/task-list

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/task-list /go/bin/task-list
ENTRYPOINT ["/go/bin/task-list"]
CMD ["server"]
