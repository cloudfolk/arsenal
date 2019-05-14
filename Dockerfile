# First stage container
FROM golang:1.12-alpine as builder

ARG APP
RUN apk add --no-cache make

ARG PKG
ADD . $GOPATH/src/${PKG}
RUN cd $GOPATH/src/${PKG} && make ${APP} && mkdir -p /build/bin && mv build/bin/* /build/bin

# Second stage container
FROM alpine:3.9

RUN apk add --no-cache ca-certificates
COPY --from=builder /build/bin/* /usr/local/bin/

ARG USER=cloudfolk
RUN adduser -D ${USER}
USER ${USER}

# Define your entrypoint or command
# ENTRYPOINT [""]
