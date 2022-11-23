FROM golang:1.19.3-alpine3.16 as builder
LABEL maintainer="mgdevstack" \
    vendor="Zettabytes" \
    owner="zbio" \
    type="example"
ADD main.go  .
RUN GO111MODULE=off CGO_ENABLED=0 GOOS=linux go build -o app

FROM alpine:3.16.2
LABEL maintainer="mgdevstack" \
    vendor="Zettabytes" \
    owner="zbio"
COPY --from=builder /go/app /app
USER nobody
ENTRYPOINT ["/app"]
