FROM golang:1.17 as build

WORKDIR /go/src/github.com/microsoft/azure-iothub-exporter

COPY ./go.mod /go/src/github.com/microsoft/azure-iothub-exporter
COPY ./go.sum /go/src/github.com/microsoft/azure-iothub-exporter
RUN go mod download

COPY ./ /go/src/github.com/microsoft/azure-iothub-exporter
RUN CGO_ENABLED=0 go build

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/github.com/microsoft/azure-iothub-exporter/azure-iothub-exporter /
USER 1000:1000
EXPOSE 8080
ENTRYPOINT ["/azure-iothub-exporter"]
