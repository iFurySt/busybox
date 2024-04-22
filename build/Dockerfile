FROM golang:1.22.2 as go-builder

WORKDIR /go/src/busybox

COPY ./ ./

RUN make all

FROM busybox:1.36.1-glibc as runtime

WORKDIR /app

COPY --from=go-builder /go/src/busybox/bin/busybox /app/busybox

CMD ["/app/busybox"]
