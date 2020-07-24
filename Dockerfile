FROM golang:1.14.4-alpine3.12

RUN addgroup -S sample && adduser -S sample -G sample

RUN mkdir "/src"
WORKDIR "/src"

COPY main main

CMD "./main"