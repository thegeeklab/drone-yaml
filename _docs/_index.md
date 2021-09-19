---
title: drone-yaml
---

[![Build Status](https://img.shields.io/drone/build/thegeeklab/drone-yaml?logo=drone&server=https%3A%2F%2Fdrone.thegeeklab.de)](https://drone.thegeeklab.de/thegeeklab/drone-yaml)
[![Docker Hub](https://img.shields.io/badge/dockerhub-latest-blue.svg?logo=docker&logoColor=white)](https://hub.docker.com/r/thegeeklab/drone-yaml)
[![Quay.io](https://img.shields.io/badge/quay-latest-blue.svg?logo=docker&logoColor=white)](https://quay.io/repository/thegeeklab/drone-yaml)
[![GitHub contributors](https://img.shields.io/github/contributors/thegeeklab/drone-yaml)](https://github.com/thegeeklab/drone-yaml/graphs/contributors)
[![Source: GitHub](https://img.shields.io/badge/source-github-blue.svg?logo=github&logoColor=white)](https://github.com/thegeeklab/drone-yaml)
[![License: MIT](https://img.shields.io/github/license/thegeeklab/drone-yaml)](https://github.com/thegeeklab/drone-yaml/blob/main/LICENSE)

Custom linter and formatter for the [Drone](https://github.com/drone/drone) YAML configuration file format.

<!-- prettier-ignore-start -->
<!-- spellchecker-disable -->
{{< toc >}}
<!-- spellchecker-enable -->
<!-- prettier-ignore-end -->

## Build

Build the binary with the following command:

```Shell
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

go build -v -a -tags netgo -o release/linux/amd64/drone-yaml
```

Build the Docker image with the following command:

```Shell
docker build --file docker/Dockerfile.amd64 --tag thegeeklab/drone-yaml .
```

## Usage

{{< hint warning >}}
**Note**\
Be aware that the tool only supports configuration files for the Drone Docker runner!
{{< /hint >}}

Lint the YAML file:

```Shell
drone-yaml lint samples/simple.yml
```

Format the YAML file:

```Shell
# default is printing to stdout
drone-yaml fmt samples/simple.yml

# optionally update the formatted file in place
drone-yaml fmt samples/simple.yml --save
```
