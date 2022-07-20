FROM ubuntu:focal AS release
ENV USER_UID=1001 USER_NAME=opslevel
ENTRYPOINT ["/usr/local/bin/opslevel"]
WORKDIR /app
RUN apt-get update && \
    apt-get install -y curl jq && \
    apt-get purge && apt-get clean && apt-get autoclean
COPY opslevel /usr/local/bin