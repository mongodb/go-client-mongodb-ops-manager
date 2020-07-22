# Contributing to MongoDB Ops Manager Go Client

Thanks for your interest in contributing to this project, 
this document describes some guidelines necessary to participate in the community.

## Feature Requests

We welcome any feedback or feature request, to submit yours
please head over to our [feedback page](https://feedback.mongodb.com/).

## Reporting Issues

Please create a [GitHub issue](https://github.com/mongodb/go-client-mongodb-ops-manager/issues) describing the kind of problem you're facing
with as much detail as possible, including things like operating system or anything else may be relevant to the issue.

## Submitting a Patch

Before submitting a patch to the repo please consider opening an [issue first](#reporting-issues)

### Contributor License Agreement

For patches to be accepted, contributors must sign our [CLA](https://www.mongodb.com/legal/contributor-agreement).

## Development setup

### Prerequisite Tools 
- [Git](https://git-scm.com/)
- [Go (at least Go 1.12)](https://golang.org/dl/)

### Environment
- Fork the repository.
- Clone your forked repository locally.
- We use Go Modules to manage dependencies, so you can develop outside of your `$GOPATH`.

We use [golangci-lint](https://github.com/golangci/golangci-lint) to lint our code, you can install it locally via `make setup`.

## Building and testing

The following is a short list of commands that can be run in the root of the project directory

- Run `make` see a list of available targets.
- Run `make test` to run all unit tests.
- Run `make lint` to validate against our linting rules.

We provide a git pre-commit hook to format and check the code, to install it run `make link-git-hooks` 

## Third party dependencies

We scan our dependencies for vulnerabilities and incompatible licenses using [Snyk](https://snyk.io/).
To run Snyk locally please follow their [CLI reference](https://support.snyk.io/hc/en-us/articles/360003812458-Getting-started-with-the-CLI) 
