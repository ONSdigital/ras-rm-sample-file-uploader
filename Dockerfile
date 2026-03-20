FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS build-stage

RUN mkdir "/src"
WORKDIR "/src"

COPY . .

RUN go build -v -o main
RUN chmod 755 main

FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS final-stage

RUN addgroup -S sample-group && adduser -S sample-user -G sample-group
RUN mkdir -p "/opt/sample"
RUN chown sample-user:sample-group /opt/sample

WORKDIR "/opt/sample"
COPY --from=build-stage /src/main .
RUN chmod 550 /opt/sample/main
RUN chown sample-user:sample-group /opt/sample/main

USER sample-user

CMD "./main"
