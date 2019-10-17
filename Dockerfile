FROM golang:1.12.9-alpine

RUN apk add --no-cache --update alpine-sdk

COPY . /go/src/go_starter
RUN cd /go/src/go_starter && make release-binary

FROM alpine:3.9

RUN apk add --no-cache ca-certificates openssl

COPY --from=0 /go/bin/go_starter /tmp/go_starter
RUN install -m 0755 /tmp/go_starter /usr/local/bin/go_starter && \
    rm /tmp/go_starter && \
    /usr/local/bin/go_starter --version

USER 1001:1001
ENTRYPOINT ["/usr/local/bin/go_starter"]
