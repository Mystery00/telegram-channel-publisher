FROM golang:1.21-alpine as builder
COPY . /usr/local/go/src/telegram-channel-publisher
WORKDIR /usr/local/go/src/telegram-channel-publisher
RUN GO111MODULE=on go build -o /usr/bin/telegram-channel-publisher telegram-channel-publisher

###
FROM alpine as final
RUN apk add --no-cache ca-certificates && \
        update-ca-certificates
WORKDIR /app
ENTRYPOINT ["/usr/bin/telegram-channel-publisher"]
COPY --from=builder /usr/bin/telegram-channel-publisher /usr/bin/
COPY templates /app/templates