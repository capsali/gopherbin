FROM golang:1.16-rc AS builder
LABEL stage=builder

RUN curl -sL https://deb.nodesource.com/setup_15.x > /tmp/setup_node.sh
RUN /bin/bash /tmp/setup_node.sh
RUN apt-get update && apt-get -y install git make git nodejs bzip2
RUN npm install --global yarn

ADD . /build/gopherbin
WORKDIR /build/gopherbin

RUN mkdir /tmp/go
ENV GOPATH /tmp/go

# build gopher binary
RUN make all-ui 

RUN chmod +x /tmp/go/bin/gopherbin

# creating a minimal image
FROM alpine

# Add bash
RUN apk add --no-cache bash su-exec libc6-compat gettext libintl

# Copy our binary to the image
COPY --from=builder /tmp/go/bin/gopherbin /usr/local/bin/gopherbin

# Add to path
ENV PATH="/usr/local/bin/gopherbin:${PATH}"

# Create gopherbin dir
RUN mkdir -p /templates && mkdir -p /etc/gopherbin && mkdir -p /secrets

# Copy templates and confd metadata
ADD docker/templates/ /templates/

# Copy entrypoint script
ADD docker/entrypoint.sh /entrypoint.sh

# Run entrypoint script
ENTRYPOINT ["/entrypoint.sh"]

# Run binary and expose port
CMD ["/usr/local/bin/gopherbin", "-config", "/etc/gopherbin/gopherbin-config.toml"]

EXPOSE 9997/tcp
