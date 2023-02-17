FROM registry.access.redhat.com/ubi9/ubi-minimal:9.1 AS builder
USER root
RUN microdnf install -y tar gzip make which

# install go
RUN curl -O -J https://dl.google.com/go/go1.19.6.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.19.6.linux-amd64.tar.gz
RUN ln -s /usr/local/go/bin/go /usr/local/bin/go

WORKDIR /opt/seatdemo

COPY . ./

RUN go mod vendor
RUN make binary

FROM registry.access.redhat.com/ubi9/ubi-minimal:9.1

COPY --from=builder /opt/seatdemo /opt/seatdemo/

EXPOSE 3000

USER 1001

ENTRYPOINT ["/opt/seatdemo/seatdemo"]

LABEL name="seatdemo" \
      version="$VERSION" \
      summary="test seatdemo service" \
      description="seatdemotestbuild"