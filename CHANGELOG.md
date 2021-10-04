
<a name="v0.3.1"></a>
## [v0.3.1] - 2021-10-04
### Docs
- cleanup readme to bring it more inline with kubectl-opslevel

### Feature
- set custom graphql user agent extras


<a name="v0.3.0"></a>
## [v0.3.0] - 2021-10-02

<a name="v0.3.0-beta.5"></a>
## [v0.3.0-beta.5] - 2021-10-02

<a name="v0.3.0-beta.4"></a>
## [v0.3.0-beta.4] - 2021-10-01

<a name="v0.3.0-beta.3"></a>
## [v0.3.0-beta.3] - 2021-10-01

<a name="v0.3.0-beta.2"></a>
## [v0.3.0-beta.2] - 2021-10-01

<a name="v0.3.0-beta.1"></a>
## [v0.3.0-beta.1] - 2021-10-01
### Bugfix
- fix outstanding multiline string issues with `export terraform`


<a name="v0.3.0-beta"></a>
## [v0.3.0-beta] - 2021-10-01
### Bugfix
- if multiline string from opslevel does not end with \n then have terraform treat it as an escaped string

### Docs
- add installation instructions for Deb/RPM

### Feature
- add deb/rpm package releases
- add create, get, list and delete for filters
- implement create, get, list, and delete for team, team member and team contact
- implement get & list for repository
- add ability to output list data as a json array
- implement correct output formatting for list tier, lifecycle and tools


<a name="v0.2.0-beta"></a>
## [v0.2.0-beta] - 2021-09-18
### Feature
- upgrade opslevel-go to v0.3.3
- initial pass at `export terraform` for exporting data from your account to be controlled by terraform
- add shell completion generation command
- add get and list check commands
- add list commands which differ from get commands
- add commands for rubric categories and levels
- add get and delete service commands
- add gpg signing
- add ability to query account lifecycles, tiers, teams and tools

### Refactor
- convert prefered environment variable prefix from `OL_` to `OPSLEVEL_` but still support old prefix
- use args instead of flags for rubric commands
- seperate get and list commands


<a name="v0.1.0-beta.5"></a>
## [v0.1.0-beta.5] - 2021-07-10

<a name="v0.1.0-beta.4"></a>
## [v0.1.0-beta.4] - 2021-07-10

<a name="v0.1.0-beta.3"></a>
## [v0.1.0-beta.3] - 2021-07-10

<a name="v0.1.0-beta.2"></a>
## [v0.1.0-beta.2] - 2021-07-10
### Refactor
- switch to goreleaser


<a name="v0.1.0-beta.1"></a>
## [v0.1.0-beta.1] - 2021-06-26
### Bugfix
- add yaml struct tags for consistent configfile parsing

### Docs
- rewrite examples to show both yaml and env var examples as well as switching to the full integration url rather then just ID

### Refactor
- standardize inputs across flags, env vars, and yaml.  Also include more imports to support a wider varity of source input data


<a name="v0.0.1-beta.2"></a>
## [v0.0.1-beta.2] - 2021-06-19
### Release
- use aws ecr alias
- remove github docker registry publish
- add windows binary cross compile


<a name="v0.0.1-beta.1"></a>
## v0.0.1-beta.1 - 2021-06-19
### Docs
- fix amazon ECR links to use new repo alias
- flesh out readme with usage and install instructions

### Feature
- add ability to scrape git commit info if available


[v0.3.1]: https://github.com/OpsLevel/cli/compare/v0.3.0...v0.3.1
[v0.3.0]: https://github.com/OpsLevel/cli/compare/v0.3.0-beta.5...v0.3.0
[v0.3.0-beta.5]: https://github.com/OpsLevel/cli/compare/v0.3.0-beta.4...v0.3.0-beta.5
[v0.3.0-beta.4]: https://github.com/OpsLevel/cli/compare/v0.3.0-beta.3...v0.3.0-beta.4
[v0.3.0-beta.3]: https://github.com/OpsLevel/cli/compare/v0.3.0-beta.2...v0.3.0-beta.3
[v0.3.0-beta.2]: https://github.com/OpsLevel/cli/compare/v0.3.0-beta.1...v0.3.0-beta.2
[v0.3.0-beta.1]: https://github.com/OpsLevel/cli/compare/v0.3.0-beta...v0.3.0-beta.1
[v0.3.0-beta]: https://github.com/OpsLevel/cli/compare/v0.2.0-beta...v0.3.0-beta
[v0.2.0-beta]: https://github.com/OpsLevel/cli/compare/v0.1.0-beta.5...v0.2.0-beta
[v0.1.0-beta.5]: https://github.com/OpsLevel/cli/compare/v0.1.0-beta.4...v0.1.0-beta.5
[v0.1.0-beta.4]: https://github.com/OpsLevel/cli/compare/v0.1.0-beta.3...v0.1.0-beta.4
[v0.1.0-beta.3]: https://github.com/OpsLevel/cli/compare/v0.1.0-beta.2...v0.1.0-beta.3
[v0.1.0-beta.2]: https://github.com/OpsLevel/cli/compare/v0.1.0-beta.1...v0.1.0-beta.2
[v0.1.0-beta.1]: https://github.com/OpsLevel/cli/compare/v0.0.1-beta.2...v0.1.0-beta.1
[v0.0.1-beta.2]: https://github.com/OpsLevel/cli/compare/v0.0.1-beta.1...v0.0.1-beta.2
