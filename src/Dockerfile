FROM ubuntu:noble AS release
ENV USER_UID=1001 USER_NAME=opslevel
ENTRYPOINT ["/usr/local/bin/opslevel"]
WORKDIR /app
# Bootstrap CA certs over plain HTTP with retries off (CVE-2026-109518 needs them), then switch apt to HTTPS.
RUN apt-get -o Acquire::Retries=0 update && \
    apt-get -o Acquire::Retries=0 install -y ca-certificates && \
    sed -i 's|http://|https://|g' /etc/apt/sources.list.d/ubuntu.sources && \
    apt-get update && \
    apt-get install -y jq && \
    apt-get purge && apt-get clean && apt-get autoclean
COPY opslevel /usr/local/bin