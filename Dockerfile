
FROM golang:1.9-alpine AS build-env

RUN apk --no-cache --virtual .builddeps add git make
RUN adduser app -g ,,, -s /bin/false -HD

ADD . /go/src/github.com/else/goertzcard-notify
RUN cd /go/src/github.com/else/goertzcard-notify; \
    make test install

FROM alpine
WORKDIR /app
COPY --from=build-env /go/bin/goertzcard-notify /app
COPY --from=build-env /etc/passwd /etc/passwd
USER app
ENTRYPOINT ["/app/goertzcard-notify"]