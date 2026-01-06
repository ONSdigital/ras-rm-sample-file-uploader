FROM --platform=$BUILDPLATFORM golang:1.25-alpine

RUN addgroup -S sample-group && adduser -S sample-user -G sample-group

RUN mkdir -p "/opt/sample"
RUN chown sample-user:sample-group /opt/sample

WORKDIR "/opt/sample"

RUN if [ "$BUILDPLATFORM" = "linux/amd64" ]; then \
      cp build/linux-amd64/bin/main .; \
    elif [ "$BUILDPLATFORM" = "linux/arm64" ]; then \
      cp build/linux-arm64/bin/main .; \
    fi

RUN chmod 550 /opt/sample/main
RUN chown sample-user:sample-group /opt/sample/main

USER sample-user

CMD "./main"
