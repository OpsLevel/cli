FROM golang:1.16 AS builder
ARG VERSION=development
LABEL stage=builder
WORKDIR /workspace
COPY ./src/go.mod .
COPY ./src/go.sum .
RUN go mod download
COPY ./src .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./opslevel -ldflags="-X 'github.com/opslevel/cli/cmd.version=${VERSION}'"


FROM ubuntu:focal AS release
ENV USER_UID=1001 USER_NAME=opslevel
ENTRYPOINT ["/usr/local/bin/opslevel"]
WORKDIR /app
RUN apt-get update && \
    apt-get install -y curl && \
    apt-get purge && apt-get clean && apt-get autoclean && \
    curl -o /usr/local/bin/jq http://stedolan.github.io/jq/download/linux64/jq && \
    chmod +x /usr/local/bin/jq
COPY --from=builder /workspace/opslevel /usr/local/bin/

