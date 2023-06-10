# sid-slim
FROM debian@sha256:10e2808dfaf6451c0ba0808fc70b2a541d1d6d0b91c8aba10b680e9778ecb4bb

ARG GO_VERSION="1.19.10"

RUN \
	apt update \
	&& apt install --no-install-recommends -y -q \
		ca-certificates \
		gcc libc-dev \
		make \
		curl git \
		libalpm-dev libarchive-dev \
	&& apt-get -y autoremove \
	&& apt-get clean \
	&& rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN rm -rf /usr/local/go \
	&& curl -s -L -O "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" \
	&& tar -C /usr/local -zxf "go${GO_VERSION}.linux-amd64.tar.gz" \
	&& rm "go${GO_VERSION}.linux-amd64.tar.gz"

ENV PATH=$PATH:/usr/local/go/bin

