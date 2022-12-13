<p align="center">
    <a href="https://github.com/OpsLevel/cli/blob/main/LICENSE" alt="License">
        <img src="https://img.shields.io/github/license/OpsLevel/cli.svg" /></a>
    <a href="https://goreportcard.com/report/github.com/OpsLevel/cli" alt="Go Report Card">
        <img src="https://goreportcard.com/badge/github.com/OpsLevel/cli" /></a>
    <a href="https://GitHub.com/OpsLevel/cli/releases/" alt="Release">
        <img src="https://img.shields.io/github/v/release/OpsLevel/cli" /></a>  
    <a href="https://masterminds.github.io/stability/experimental.html" alt="Stability: Experimental">
        <img src="https://masterminds.github.io/stability/experimental.svg" /></a>  
    <a href="https://github.com/OpsLevel/cli/graphs/contributors" alt="Contributors">
        <img src="https://img.shields.io/github/contributors/OpsLevel/cli" /></a>
    <a href="https://github.com/OpsLevel/cli/pulse" alt="Activity">
        <img src="https://img.shields.io/github/commit-activity/m/OpsLevel/cli" /></a>
    <a href="https://github.com/OpsLevel/cli/releases" alt="Downloads">
        <img src="https://img.shields.io/github/downloads/OpsLevel/cli/total" /></a>
</p>

[![Overall](https://img.shields.io/endpoint?style=flat&url=https%3A%2F%2Fapp.opslevel.com%2Fapi%2Fservice_level%2FEaWapOq9VQj5FvymQEgCPNJcbF-TOibHn89Arw7d_OY)](https://app.opslevel.com/services/opslevel_cli/maturity-report)

# The CLI for interacting with [OpsLevel](https://www.opslevel.com/)




### Prerequisite

- [jq](https://stedolan.github.io/jq/download/)
- [OpsLevel API Token](https://app.opslevel.com/api_tokens)
  - Generate token by clicking `Create API Token` and providing a description
  - Export the API Token for cli access:
    ```sh
    > export OPSLEVEL_API_TOKEN=<api_token>
    ```

### Installation

#### MacOS

```sh
brew install opslevel/tap/cli
```

<!--
#### Deb

```sh
sudo apt-get install apt-transport-https
wget -qO - https://opslevel.github.io/cli-repo/deb/public.key | sudo apt-key add -
echo deb https://opslevel.github.io/cli-repo/deb [CODE_NAME] main | sudo tee -a /etc/apt/sources.list
sudo apt-get update
sudo apt-get install opslevel
```

#### RPM

```sh
cat << EOF > /etc/yum.repos.d/opslevel.repo
[opslevel]
name=opslevel cli repository
baseurl=https://opslevel.github.io/cli-repo/rpm/releases/$releasever/$basearch/
gpgcheck=0
enabled=1
EOF
sudo yum -y update
sudo yum -y install opslevel
```
-->

#### Docker

The docker container is hosted on [AWS Public ECR](https://gallery.ecr.aws/opslevel/cli)

### Quickstart

```sh
# Create
> opslevel create category Chaos
Z2lkOi8vb3BzbGV2ZWwvQ2F0ZWdvcnkvOTY8
# Get
> opslevel get category Z2lkOi8vb3BzbGV2ZWwvQ2F0ZWdvcnkvOTY8
{
  "id": "Z2lkOi8vb3BzbGV2ZWwvQ2F0ZWdvcnkvOTY8",
  "Name": "Chaos"
}
# List
> opslevel list category
NAME            ID                                    
Performance     Z2lkOi8vb3BzbGV2ZWwvQ2F0ZWdvcnkvOTY1  
Infrastructure  Z2lkOi8vb3BzbGV2ZWwvQ2F0ZWdvcnkvOTY2  
Observability   Z2lkOi8vb3BzbGV2ZWwvQ2F0ZWdvcnkvOTY3  
Reliability     Z2lkOi8vb3BzbGV2ZWwvQ2F0ZWdvcnkvOTY4  
Scalability     Z2lkOi8vb3BzbGV2ZWwvQ2F0ZWdvcnkvOTY5  
Security        Z2lkOi8vb3BzbGV2ZWwvQ2F0ZWdvcnkvOTY6  
Quality         Z2lkOi8vb3BzbGV2ZWwvQ2F0ZWdvcnkvOTY7  
Chaos           Z2lkOi8vb3BzbGV2ZWwvQ2F0ZWdvcnkvOTY8  
# Delete
> opslevel delete category Z2lkOi8vb3BzbGV2ZWwvQ2F0ZWdvcnkvOTY8
```

<blockquote>This tool is still in beta.</blockquote>

### Enable shell autocompletion

We have the ability to generate autocompletion scripts for the shell's `bash`, `zsh`, `fish` and `powershell`.  To generate 
the completion script for MacOS zsh:

```sh
opslevel completion zsh > /usr/local/share/zsh/site-functions/_opslevel
```

Make sure you have `zsh` completion turned on by having the following as one of the first few lines in your `.zshrc` file

```sh
echo "autoload -U compinit; compinit" >> ~/.zshrc
```

<!--
### JSON-Schema
TODO
-->

## Troubleshooting

### List all my tier 1 services

```sh
> opslevel list services -o json | jq '[.[] | if .tier.Alias == "tier_1" then {(.name) : (.tier.Alias)} else empty end]' 
[
  {
    "Catalog Service": "tier_1"
  },
  {
    "Shopping Cart Service": "tier_1"
  },
  {
    "Website Aggregator": "tier_1"
  }
]
```
