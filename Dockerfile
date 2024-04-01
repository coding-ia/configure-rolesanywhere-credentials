FROM golang:1.21-alpine AS builder

ENV GO111MODULE=on \
  CGO_ENABLED=1 \
  GOOS=linux \
  GOARCH=amd64

RUN apk update && apk upgrade
RUN apk add upx gcc musl-dev

WORKDIR /src
COPY . .

RUN go build \
  -ldflags "-s -w -extldflags '-static'" \
  -o /bin/ra_action \
  . \
  && strip /bin/ra_action \
  && upx -q -9 /bin/ra_action

RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc_passwd

FROM scratch

COPY --from=builder /etc_passwd /etc/passwd
COPY --from=builder --chown=65534:0 /bin/ra_action /ra_action

USER nobody
ENTRYPOINT ["/ra_action"]