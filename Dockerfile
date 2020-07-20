FROM golang:1.14 as builder

WORKDIR /go/src/app
COPY . .
ENV CGO_ENABLED=0
RUN make getmodules
RUN make buildlocal
RUN pwd
RUN ls -alh bin/

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/app/bin/cogo /cogo

ENTRYPOINT ["/cogo"]
