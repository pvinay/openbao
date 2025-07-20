FROM debian:bookworm AS default

ARG BIN_NAME=bao
ARG NAME=openbao
ARG PRODUCT_VERSION=0.0.1
ARG PRODUCT_REVISION=1
ARG VERSION=0.0.1

LABEL name="OpenBao" \
      maintainer="OpenBao <openbao@lists.openssf.org>" \
      vendor="OpenBao" \
      version=${PRODUCT_VERSION} \
      release=${PRODUCT_REVISION} \
      revision=${PRODUCT_REVISION} \
      summary="OpenBao is a tool for securely accessing secrets." \
      description="OpenBao is a tool for securely accessing secrets. A secret is anything that you want to tightly control access to, such as API keys, passwords, certificates, and more. OpenBao provides a unified interface to any secret, while providing tight access control and recording a detailed audit log."

COPY LICENSE /licenses/mozilla.txt

# Set ARGs as ENV so they can be used in ENTRYPOINT/CMD
ENV NAME=$NAME
ENV VERSION=$VERSION

# Create a non-root user to run the software.
RUN groupadd ${NAME} && useradd -r -g ${NAME} ${NAME}

ARG EXTRA_PACKAGES
RUN apt-get update && \
    apt-get install -y --no-install-recommends wget libcap2-bin dumb-init tzdata ${EXTRA_PACKAGES} && \
    rm -rf /var/lib/apt/lists/*

# golang 1.23
WORKDIR /tmp

RUN wget --no-check-certificate https://go.dev/dl/go1.24.5.linux-amd64.tar.gz 
RUN rm -rf /usr/local/go
RUN tar -xvzf go1.24.5.linux-amd64.tar.gz -C /usr/local/
#RUN cp go/bin/go /usr/local/bin/go 

RUN echo "export GO_ROOT=/usr/local/go" >> ~/.bashrc
RUN echo "export PATH=/usr/local/go/bin:$PATH" >> ~/.bashrc

RUN cp /usr/local/go/bin/go /usr/local/bin/go

# Reset working dir to repo base
WORKDIR /root/repo/openbao

COPY . .
