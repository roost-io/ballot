FROM golang:alpine as builder
LABEL maintainer="mgdevstack" \
    vendor="Zettabytes" \
    owner="zbio" \
    type="example"
ADD main.go  .
ADD main_test.go  .
RUN GO111MODULE=off CGO_ENABLED=0 GOOS=linux go build -o ballot
RUN GO111MODULE=off CGO_ENABLED=0 GOOS=linux go test -c -o ballot.test

# To execute unit test in container
FROM alpine as ballottest
LABEL maintainer="mgdevstack" \
    vendor="Zettabytes" \
    owner="zbio"
COPY --from=builder /go/ballot.test /usr/local/bin/ballot.test
USER nobody
CMD ["ballot.test", "-test.v"]

## Main ballot application container
FROM alpine
LABEL maintainer="mgdevstack" \
    vendor="Zettabytes" \
    owner="zbio"
COPY --from=builder /go/ballot /usr/local/bin/ballot
USER nobody
CMD [ "ballot" ] 