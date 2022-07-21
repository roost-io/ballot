FROM golang:1.19rc2-alpine3.16 as builder
LABEL maintainer="mgdevstack" \
    vendor="Zettabytes" \
    owner="zbio" \
    type="example"
ADD main.go  .
RUN GO111MODULE=off CGO_ENABLED=0 GOOS=linux go build -o app

FROM alpine:3.16.1
LABEL maintainer="mgdevstack" \
    vendor="Zettabytes" \
    owner="zbio"
COPY --from=builder /go/app /app
USER nobody
ENTRYPOINT ["/app"]
