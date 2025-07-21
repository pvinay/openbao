FROM debian:bookworm AS default


ARG EXTRA_PACKAGES
RUN apt-get update && \
    apt-get install -y --no-install-recommends git wget libcap2-bin dumb-init tzdata ${EXTRA_PACKAGES} && \
    rm -rf /var/lib/apt/lists/*

# golang 1.23
WORKDIR /tmp

RUN wget --no-check-certificate https://go.dev/dl/go1.24.5.linux-amd64.tar.gz 
RUN rm -rf /usr/local/go
RUN tar -xzf go1.24.5.linux-amd64.tar.gz -C /usr/local/
RUN cp /usr/local/go/bin/go /usr/local/bin/go

ENV GO_ROOT=/usr/local/go
ENV PATH="/usr/local/go/bin:$PATH"

WORKDIR /root/repo/openbao

COPY . .

RUN git config --global core.autocrlf input
RUN git status --porcelain


