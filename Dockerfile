FROM buildpack-deps:jessie-scm
MAINTAINER John Ward "jward@launchpadcentral.com"

RUN apt-get update && apt-get install -y --no-install-recommends \
		g++ \
		gcc \
		libc6-dev \
		make \
		pkg-config \
    bison \
	&& rm -rf /var/lib/apt/lists/*

SHELL ["/bin/bash", "-c"]

ENV GO_VERSION 1.12

RUN curl -s -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer | bash

RUN . /root/.gvm/scripts/gvm && \
      gvm install go1.4 -B && \
      gvm use go1.4 && \
      gvm install go1.12

ENV WATCHER_VERSION 0.2.4

ADD https://github.com/canthefason/go-watcher/releases/download/v${WATCHER_VERSION}/watcher-${WATCHER_VERSION}-linux-amd64 /root/.gvm/bin/watcher

RUN chmod +x /root/.gvm/bin/watcher

ENV GOPATH /go

WORKDIR /go/src

VOLUME /go/src
ADD entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]
CMD ["watcher"]
