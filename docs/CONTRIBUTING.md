# Want to contribute?

Great! We welcome contributions of all kinds, big or small! This includes bug reports, code fixes, documentation
improvements, and code examples.

Before you dive in, please take a moment to read through this guide.

# Reporting issue

We use [GitHub](https://github.com/5GSEC/SentryFlow) to manage the issues. Please open
a [new issue](https://github.com/5GSEC/SentryFlow/issues/new/choose) directly there.

# Getting Started

## Setting Up Your Environment

- Head over to [GitHub](https://github.com/5GSEC/SentryFlow) and fork the 5GSec SentryFlow repository.
- Clone your forked repository onto your local machine.
  ```shell
  git clone git@github.com:<your-username>/SentryFlow.git
  ```

## Install development tools

You'll need these tools for a smooth development experience:

- [Make](https://www.gnu.org/software/make/#download)
- [Go](https://go.dev/doc/install) SDK, version 1.23 or later
- Go IDE ([Goland](https://www.jetbrains.com/go/) / [VS Code](https://code.visualstudio.com/download))
- Container tools ([Docker](https://www.docker.com/) / [Podman](https://podman.io/))
- [Kubernetes cluster](https://kubernetes.io/docs/setup/) running version 1.28 or later.
- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) version 1.28 or later.

# Contributing Code

## Building Locally

- Install development tools as [mentioned above](#setting-up-your-environment).

- Build SentryFlow using:
  ```shell
  cd sentryflow
  make build
  ```

## Understanding the Project

Before contributing to any Open Source project, it's important to have basic understanding of what the project is about.
It is advised to try out the project as an end user.

## Project Structure

These are general guidelines on how to organize source code in this repository.

```
github.com/5GSEC/SentryFlow

├── client                                    -> Log client code.
├── deployments                               -> Manifests or Helm charts for deployment on Kubernetes.
├── docs                                      -> All Documentation.
│   └── receivers                                 -> Receiver specifc integration documentaion.
│       ├── other 
│       │   ├── ingress-controller
│       │   │   └── nginx-inc
│       │   └── web-server
│       │       └── nginx
│       └── service-mesh
│           └── istio
├── filter                                    -> Receivers specific filters/modules to observe API calls from receivers.
├── protobuf
│   ├── golang                                -> Generated protobuf Go code.
│   ├── python                                -> Generated protobuf Python code.
├── scripts
├── sentryflow
│   ├── cmd                                   -> Code for the actual binary.
│   ├── config
│   │   └── default.yaml                     -> Default configuration file.
│   ├── go.mod                               -> Go module file to track dependencies.
│   └── pkg                                  -> pkg is a collection of utility packages used by the components without being specific to its internals.
│       ├── config                           -> Configuration initialization code.
│       ├── core                             -> SentryFlow core initialization code.
│       ├── exporter                         -> Exporter code.
│       ├── k8s                              -> Kubernetes client code.
│       ├── receiver                         -> Receiver code.
│       │   ├── receiver.go                     -> All receivers initialization code.
│       │   └── svcmesh                         -> ServiceMesh receivers code.
│       │   └── other                           -> Other receivers code.
│       └── util                            -> Utilities.
```

## Imports grouping

This project follows the following pattern for grouping imports in Go files:

* imports from standard library.
* imports from other projects.
* imports from `sentryflow` project.

For example:

```go
import (
  "context"
  "fmt"
  
  "k8s.io/apimachinery/pkg/runtime"
  "sigs.k8s.io/controller-runtime/pkg/client"
  
  "github.com/5GSEC/SentryFlow/pkg/config"
  "github.com/5GSEC/SentryFlow/pkg/receiver"
  "github.com/5GSEC/SentryFlow/pkg/util"
)
```

## Pull Requests and Code Reviews

We use GitHub [pull requests](https://github.com/5GSEC/SentryFlow/pulls) for code contributions. All submissions,
including
those from project members, require review before merging.
We typically aim for two approvals per pull request, with reviews happening within a week or two.
Feel free to ping reviewers if you haven't received feedback within that timeframe.

### Commit Messages

We follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification for clear and
consistent commit messages.

Please make sure you have added the **Signed-off-by:** footer in your git commit. In order to do it automatically, use
the **--signoff** flag:

```shell
git commit --signoff
```

With this command, `git` would automatically add a footer by reading your name and email from your `.gitconfig` file.

### Merging PRs

**For maintainers:** Before merging a PR make sure the title is descriptive and follows
a [good commit message](https://www.conventionalcommits.org/en/v1.0.0/).

Merge the PR by using `Squash and merge` option on GitHub. Avoid creating merge commits. After the merge make sure
referenced issues were closed.

# Testing and Documentation

Tests and documentation are not optional, make sure your pull requests include:

- Tests that verify your changes and don't break existing functionality.
- Updated [documentation](../docs) reflecting your code changes.
- Reference information and any other relevant details.

## Commands to run tests

- Unit tests:
  ```shell
  make tests
  ```

- Integration tests:
  ```shell
  make integration-test
  ```

- End-to-end tests:
  ```shell
  make e2e-test
  ```

