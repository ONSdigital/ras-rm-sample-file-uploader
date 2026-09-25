FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS build-stage
ARG TARGETOS
ARG TARGETARCH

RUN mkdir "/src"
WORKDIR "/src"

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -v -o main
RUN chmod 755 main

FROM alpine:3.24 AS final-stage

RUN apk add --no-cache ca-certificates

RUN addgroup -S sample-group && adduser -S sample-user -G sample-group
RUN mkdir -p "/opt/sample"
RUN chown sample-user:sample-group /opt/sample

WORKDIR "/opt/sample"
COPY --from=build-stage /src/main .
RUN chmod 550 /opt/sample/main
RUN chown sample-user:sample-group /opt/sample/main

USER sample-user

CMD "./main"
