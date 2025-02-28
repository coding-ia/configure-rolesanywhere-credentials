FROM golang:1.24-alpine AS builder

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

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/ra_action /ra_action

ENTRYPOINT ["/ra_action"]