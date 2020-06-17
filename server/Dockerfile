FROM golang:1.14-alpine AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -mod=vendor -o /go/bin/task-list

FROM scratch
COPY --from=builder /go/bin/task-list /go/bin/task-list
ENTRYPOINT ["/go/bin/task-list"]
CMD ["server"]
