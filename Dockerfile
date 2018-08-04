FROM golang:1.10-alpine AS build
LABEL maintainer betchi

RUN apk add --update --no-cache alpine-sdk bash python
WORKDIR /root
RUN git clone -b v0.11.4 https://github.com/edenhill/librdkafka.git
WORKDIR /root/librdkafka
RUN ./configure
RUN make
RUN make install

COPY . /go/src/github.com/swagchat/rtm-api
WORKDIR /go/src/github.com/swagchat/rtm-api/
RUN go build -o rtm-api


FROM alpine:3.7
LABEL maintainer betchi

RUN apk --no-cache --update upgrade \
  && apk add --update --no-cache tzdata ca-certificates \
  && update-ca-certificates --fresh

RUN mkdir -p /app
COPY --from=build /go/src/github.com/swagchat/rtm-api/rtm-api /app/rtm-api
COPY --from=build /usr/local/lib/librdkafka.a /usr/local/lib/librdkafka.a
COPY --from=build /usr/local/lib/librdkafka.so /usr/local/lib/librdkafka.so
COPY --from=build /usr/local/lib/librdkafka.so.1 /usr/local/lib/librdkafka.so.1
COPY --from=build /usr/local/include/librdkafka /usr/local/include/librdkafka

STOPSIGNAL SIGTERM

EXPOSE 8102
CMD /bin/sh -c "cd /app && ./rtm-api"
