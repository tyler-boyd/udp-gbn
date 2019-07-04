FROM ubuntu:18.04

RUN \
  apt-get -qq update && \
  apt-get -qqy install \
    curl vim apt-utils wget tar make

RUN wget -q https://dl.google.com/go/go1.12.6.linux-amd64.tar.gz
RUN tar -xf go1.12.6.linux-amd64.tar.gz
RUN mv go /usr/local

ENV GOROOT /usr/local/go
ENV GOPATH /go

ENV PATH $GOPATH/bin:$GOROOT/bin:$PATH


WORKDIR /go/src/github.com/tyler-boyd/udp-gbn

COPY . /go/src/github.com/tyler-boyd/udp-gbn


ENTRYPOINT [ "/bin/bash" ]
