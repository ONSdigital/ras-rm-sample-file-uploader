FROM golang:1.19.3-alpine3.16

RUN addgroup -S sample-group && adduser -S sample-user -G sample-group

RUN mkdir -p "/opt/sample"
RUN chown sample-user:sample-group /opt/sample

WORKDIR "/opt/sample"

COPY main .

RUN chmod 550 /opt/sample/main
RUN chown sample-user:sample-group /opt/sample/main

USER sample-user

CMD "./main"