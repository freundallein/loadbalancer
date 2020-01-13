FROM golang:alpine AS intermediate

RUN apk update && \
    apk add --no-cache git make

RUN adduser -D -g '' lb

WORKDIR $GOPATH/src/

COPY . .

RUN go mod download
RUN go mod verify
RUN make build

FROM scratch

ENV TEST_ADDR=
ENV PORT=8000
ENV TEST_STALE_TIMEOUT=60

COPY --from=intermediate /go/src/bin/go-lb /go/bin/go-lb
COPY --from=intermediate /etc/passwd /etc/passwd

USER lb

WORKDIR /go/bin

EXPOSE $PORT

CMD ["/go/bin/go-lb"]