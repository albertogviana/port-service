FROM golang:1.17 AS builder

ADD . /go/src/github.com/albertogviana/port-service
WORKDIR /go/src/github.com/albertogviana/port-service
RUN make build-port-service-linux


FROM scratch

COPY --from=builder /go/src/github.com/albertogviana/port-service/cmd/port-service /port-service

WORKDIR /import
# nobody user
USER 65534

ENTRYPOINT [ "/port-service" ]
