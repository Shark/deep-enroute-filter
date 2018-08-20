FROM golang:1.10-stretch

RUN apt-get update \
 && apt-get install -y --no-install-recommends iptables libnetfilter-queue-dev \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* \
 && curl -sSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 \
 && chmod +x /usr/local/bin/dep \
 && mkdir -p /go/src/gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter

WORKDIR /go/src/gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter

COPY Gopkg.lock Gopkg.toml /go/src/gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/
RUN dep ensure -vendor-only

RUN cd vendor/github.com/zubairhamed/canopus \
 && git clone https://github.com/openssl/openssl.git \
 && cd openssl \
 && git checkout b9b5181 \
 && ./config && make -j4

COPY . .
RUN go build .
CMD ["/go/src/gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/filter"]
