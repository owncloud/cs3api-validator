FROM amd64/alpine:3.14

RUN apk update && \
	apk upgrade && \
	apk add ca-certificates mailcap tree attr && \
	rm -rf /var/cache/apk/* && \
	echo 'hosts: files dns' >| /etc/nsswitch.conf

LABEL maintainer="ownCloud GmbH <engineering@owncloud.com>" \
  org.opencontainers.image.title="CS3Api Validator" \
  org.opencontainers.image.vendor="ownCloud GmbH" \
  org.opencontainers.image.authors="ownCloud GmbH" \
  org.opencontainers.image.description="CS3Api Validator E2E testing tool for the CS3Apis" \
  org.opencontainers.image.licenses="Apache-2.0" \
  org.opencontainers.image.documentation="https://github.com/owncloud/cs3api-validator" \
  org.opencontainers.image.url="https://hub.docker.com/r/owncloud/cs3api-validator" \
  org.opencontainers.image.source="https://github.com/owncloud/cs3api-validator"

RUN mkdir -p /var/lib/cs3api-validator && \
    chmod -R 777 /var/lib/cs3api-validator

WORKDIR /var/lib/cs3api-validator

# the validator needs the path to the features file as a first param
ENTRYPOINT ["/usr/bin/cs3api-validator", "."]
CMD []

COPY features/* /var/lib/cs3api-validator
COPY cs3api-validator.test /usr/bin/cs3api-validator
