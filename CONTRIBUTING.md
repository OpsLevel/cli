# Contributing to opslevel-cli

👋 Welcome, and thank you for your interest in contributing to `opslevel-cli`! This guide will help you ramp up, propose changes, develop locally, and contribute code effectively.

---

## Table of Contents

1. [About the CLI](#about-the-cli)
2. [Getting Started](#getting-started)
3. [Development Workflow](#development-workflow)
4. [Proposing and Submitting Changes](#proposing-and-submitting-changes)
5. [Command Architecture & Style](#command-architecture--style)
6. [Testing & Tooling](#testing--tooling)
7. [Release Process](#release-process)

---

## About the CLI

The `opslevel-cli` is a command-line interface for interacting with the [OpsLevel](https://www.opslevel.com) API. It helps engineers automate, inspect, and manage their service catalog, ownership, checks, and more.

### Architecture

**[Cobra](https://github.com/spf13/cobra)** is the library we use to define our CLI. Commands are defined as Go structs with handlers, descriptions, and flags.
- **[Viper](https://github.com/spf13/viper)** handles flag parsing, environment variables, and configuration files.
- Modular command files live under `/cmd`, grouped by functionality (e.g., services, checks, etc.).
- Commands are registered to `rootCmd` via `init()` functions.
- 80% of our functionality is provided by opslevel-go and the purpose of this CLI is just to marshal data between the user and opslevel-go in a UX friendly way.
- Most commands follow the standard CRUD pattern `opslevel create ...`, `opslevel get ...`, `opslevel list ...`, `opslevel update ...`, `opslevel delete ...`, etc.
- We have an `opslevel beta` subcommand for experimental commands that are subject to removal.

---

## Getting Started

### Clone the Repo

If you're an external contributor fork the repo, then:

```bash
git clone https://github.com/YOUR_USERNAME/cli.git
cd cli
git checkout -b my-feature
```

If you're part of the OpsLevel org:

```bash
git clone git@github.com:OpsLevel/cli.git
cd cli
git checkout -b my-feature
```

> You might need to be given permissions to the repo, please reach out to team platform.

### Prerequisites

- [Task](https://taskfile.dev)
- An [OpsLevel API Token](https://app.opslevel.com/api_tokens)

You can use `task setup` to run an idempotent one-time setup.

Set your API token:

```sh
export OPSLEVEL_API_TOKEN=your_token_here
```

If you need to target a different environment, set the `OPSLEVEL_API_URL` environment variable:

```sh
export OPSLEVEL_API_URL=https://self-hosted.opslevel.dev/
```

---

## Development Workflow

### Run the CLI Locally

```sh
cd ./src
go run main.go --help
```

This is the easiest way to test your changes live.

### Using a commit from an `opslevel-go` Branch

Sometimes you will need to make CLI changes in tandem with changes to [`opslevel-go`](https://github.com/OpsLevel/opslevel-go). 
To test changes from a local branch of [`opslevel-go`](https://github.com/OpsLevel/opslevel-go):

```sh
# Set up workspace using task
task workspace

# Switch to your feature branch in the submodule
git -C ./src/submodules/opslevel-go checkout --track origin/my-feature-branch
```

All CLI calls will now use your local `opslevel-go` code checked out into the submodule at `./src/submodules/opslevel-go`.
This way you can effectively work on both the CLI and `opslevel-go` in parallel if needed.

## Testing & Tooling

- Use `task test` to run tests locally
- Use `task lint` to check for code quality issues locally
- Use `task fix` to fix formatting, linting, go.mod, and update submodule all in one go

Our CI pipeline will run `task ci` which can also be run locally to debug any issue that might only arise in CI.

---

## Proposing and Submitting Changes

### 1. Open an Issue

If you're fixing a bug or adding a feature, please [open an issue](https://github.com/OpsLevel/cli/issues) first. This helps us track discussions and ensures your work aligns with project goals.

> **Note:** All PRs must be tied to a GitHub issue.

Look for issues marked `good first issue` if you're new.

### 2. Submit a Pull Request

- Push your changes to your fork or feature branch
- Run `changie new` to create a changelog entry
- Open a PR to `main` or the relevant feature branch
- Our GitHub Actions pipeline will run tests and lint checks
- A maintainer will review your PR and request changes if needed

Once approved and merged, your change will be included in the next release.

---

## Command Architecture & Style

- Each CLI command lives in its own file under `/cmd`
- Commands should:
    - Register themselves via `init()`
    - Use `RunE` instead of `Run` for error handling
- Flags and environment variables should be registered with Viper for consistency
- Prefer readable, minimalistic command logic — delegate heavy logic to opslevel-go unless it's not possible

### Minimal Command Example

The following shows the minimum amount of code needed to create a command.  There are a plethora of commands already registered to the root command, so this is just an example of how to create a new command.
Please take a look at the existing commands to get a feel for how they work and whats possible.

```go
var exampleCmd = &cobra.Command{
    Use:   "example",
    Short: "Hello World Command",
    Long:  "Hello World Command to show how an example command works",
    RunE: func(cmd *cobra.Command, args []string) error {
        log.Info().Msg("Hello World!")
        return nil
    },
}

func init() {
    rootCmd.AddCommand(exampleCmd)
}
```

---

## Release Process

- All customer facing changes must have a [Changie](https://changie.dev) changelog entry
    - Run: `changie new`
    - Follow prompts to categorize your change
    - This generates a YAML file in `.changes/` that must be committed with your PR

- CI/CD (GitHub Actions) runs lint and tests automatically on pull requests and the main branch
- Maintainers will merge once approved
- Your contribution will be included in the next versioned release (triggered by a maintainer)

---

Happy hacking 🎉 and thank you for helping improve the `opslevel-cli`!
