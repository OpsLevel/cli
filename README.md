<p align="center">
    <a href="https://github.com/OpsLevel/cli/blob/main/LICENSE" alt="License">
        <img src="https://img.shields.io/github/license/OpsLevel/cli.svg" /></a>
    <a href="http://golang.org" alt="Made With Go">
        <img src="https://img.shields.io/github/go-mod/go-version/OpsLevel/cli?filename=src%2Fgo.mod" /></a>
    <a href="https://GitHub.com/OpsLevel/cli/releases/" alt="Release">
        <img src="https://img.shields.io/github/v/release/OpsLevel/cli" /></a>  
    <a href="https://GitHub.com/OpsLevel/cli/issues/" alt="Issues">
        <img src="https://img.shields.io/github/issues/OpsLevel/cli.svg" /></a>  
    <a href="https://github.com/OpsLevel/cli/graphs/contributors" alt="Contributors">
        <img src="https://img.shields.io/github/contributors/OpsLevel/cli" /></a>
    <a href="https://github.com/OpsLevel/cli/pulse" alt="Activity">
        <img src="https://img.shields.io/github/commit-activity/m/OpsLevel/cli" /></a>
    <a href="https://dependabot.com/" alt="Dependabot">
        <img src="https://badgen.net/badge/Dependabot/enabled/green?icon=dependabot" /></a>
</p>

`opslevel` is the command line tool for interacting with [OpsLevel](https://www.opslevel.com/)

### Quickstart

Follow the [installation](#installation) instructions before running the below commands

```bash
opslevel create deploy -i "https://app.opslevel.com/integrations/deploy/XXX" -s "foo"
```
OR
```bash
cat << EOF | opslevel create deploy -i "https://app.opslevel.com/integrations/deploy/XXX" -f -
service: "foo"
description: "Hello World"
environment: "Production"
deploy-number: 10
deploy-url: http://example.com
dedup-id: 123456789
deployer:
  name: glen
  email: glen@example.com
EOF
```
OR
```bash
export OL_INTEGRATION_URL="https://app.opslevel.com/integrations/deploy/XXX"
export OL_SERVICE=foo
export OL_DESCRIPTION="Hello World"
export OL_ENVIRONMENT=Production
export OL_DEPLOY_NUMBER=10
export OL_DEPLOY_URL="http://example.com"
export OL_DEDUP_ID=123456789
export OL_DEPLOYER_NAME=glen
export OL_DEPLOYER_EMAIL=glen@example.com
export OL_COMMIT_SHA=0s9df90sdf09
export OL_COMMIT_MESSAGE="Hello world"
opslevel create deploy
```

It can also be run with our public docker container

```bash
docker run -it --rm -v $(pwd):/app public.ecr.aws.com/opslevel/cli:0.0.1 create deploy -s "foo"
```

<!---
TODO: Add CLI Demo Gif
-->

<blockquote>This tool is still in beta.</blockquote>

## Installation

#### MacOS

```sh
TOOL_VERSION=$(curl -s https://api.github.com/repos/opslevel/cli/releases/latest | grep tag_name | cut -d '"' -f 4)
curl -Lo opslevel.tar.gz https://github.com/opslevel/cli/releases/download/${TOOL_VERSION}/opslevel-darwin-amd64.tar.gz
tar -xzvf opslevel.tar.gz  
rm opslevel.tar.gz
sudo mv opslevel /usr/local/bin/opslevel
```

#### Linux

```sh
TOOL_VERSION=$(curl -s https://api.github.com/repos/opslevel/cli/releases/latest | grep tag_name | cut -d '"' -f 4)
curl -Lo opslevel.tar.gz https://github.com/opslevel/cli/releases/download/${TOOL_VERSION}/opslevel-linux-amd64.tar.gz
tar -xzvf opslevel.tar.gz  
rm opslevel.tar.gz
sudo mv opslevel /usr/local/bin/opslevel
```

#### Docker

The docker container is hosted on [AWS Public ECR](https://gallery.ecr.aws/opslevel/cli)

The following downloads the container and creates a shim at `/usr/local/bin` so you can use the binary like its downloaded natively - it will just be running in a docker container. *NOTE: you may need to adjust how your kube config is mounted and set inside the container*

```
TOOL_VERSION=$(curl -s https://api.github.com/repos/opslevel/cli/releases/latest | grep tag_name | cut -d '"' -f 4)
docker pull public.ecr.aws/opslevel/cli:${TOOL_VERSION}
docker tag public.ecr.aws/opslevel/cli:${TOOL_VERSION} opslevel:latest 
cat << EOF > /usr/local/bin/opslevel
#! /bin/sh
docker run -it --rm -v \$(pwd):/app -v ${HOME}/.kube:/.kube -e KUBECONFIG=/.kube/config --network=host opslevel:latest \$@
EOF
chmod +x /usr/local/bin/opslevel
```