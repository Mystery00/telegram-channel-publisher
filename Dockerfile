FROM golang:1.21 as builder
COPY . /usr/local/go/src/telegram-channel-publisher
WORKDIR /usr/local/go/src/telegram-channel-publisher
RUN GO111MODULE=on go build -o /usr/bin/telegram-channel-publisher telegram-channel-publisher

###
FROM ubuntu:jammy as final
RUN export DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y tzdata less curl
RUN ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && dpkg-reconfigure --frontend noninteractive tzdata
RUN apt install --reinstall ca-certificates -y
WORKDIR /app
ENTRYPOINT ["/usr/bin/telegram-channel-publisher"]
COPY --from=builder /usr/bin/telegram-channel-publisher /usr/bin/
COPY templates /app/templates