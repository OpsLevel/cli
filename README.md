<p align="center">
    <a href="https://github.com/OpsLevel/cli/blob/main/LICENSE" alt="License">
        <img src="https://img.shields.io/github/license/OpsLevel/cli.svg" /></a>
    <a href="https://goreportcard.com/report/github.com/OpsLevel/cli" alt="Go Report Card">
        <img src="https://goreportcard.com/badge/github.com/OpsLevel/cli" /></a>
    <a href="https://GitHub.com/OpsLevel/cli/releases/" alt="Release">
        <img src="https://img.shields.io/github/v/release/OpsLevel/cli" /></a>  
    <a href="https://GitHub.com/OpsLevel/cli/issues/" alt="Issues">
        <img src="https://img.shields.io/github/issues/OpsLevel/cli.svg" /></a>  
    <a href="https://github.com/OpsLevel/cli/graphs/contributors" alt="Contributors">
        <img src="https://img.shields.io/github/contributors/OpsLevel/cli" /></a>
    <a href="https://github.com/OpsLevel/cli/pulse" alt="Activity">
        <img src="https://img.shields.io/github/commit-activity/m/OpsLevel/cli" /></a>
    <a href="https://github.com/OpsLevel/cli/releases" alt="Downloads">
        <img src="https://img.shields.io/github/downloads/OpsLevel/cli/total" /></a>
</p>

<p align="center">
 <a href="#quickstart">Quickstart</a> |
 <a href="#prerequisite">Prerequisite</a> |
 <a href="#installation">Installation</a> |
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
export OPSLEVEL_INTEGRATION_URL="https://app.opslevel.com/integrations/deploy/XXX"
export OPSLEVEL_SERVICE=foo
export OPSLEVEL_DESCRIPTION="Hello World"
export OPSLEVEL_ENVIRONMENT=Production
export OPSLEVEL_DEPLOY_NUMBER=10
export OPSLEVEL_DEPLOY_URL="http://example.com"
export OPSLEVEL_DEDUP_ID=123456789
export OPSLEVEL_DEPLOYER_NAME=glen
export OPSLEVEL_DEPLOYER_EMAIL=glen@example.com
export OPSLEVEL_COMMIT_SHA=0s9df90sdf09
export OPSLEVEL_COMMIT_MESSAGE="Hello world"
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

### Prerequisite

- [jq](https://stedolan.github.io/jq/download/)
- [OpsLevel API Token](https://app.opslevel.com/api_tokens)

### Installation

```sh
brew install opslevel/tap/cli
```

#### Docker

The docker container is hosted on [AWS Public ECR](https://gallery.ecr.aws/opslevel/cli)
