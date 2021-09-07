# ods-jenkinsfile-converter
[![build](https://github.com/opendevstack/ods-jenkinsfile-converter/actions/workflows/build.yml/badge.svg)](https://github.com/opendevstack/ods-jenkinsfile-converter/actions/workflows/build.yml)

## Introduction

This is a Golang CLI to convert the a Jenkinsfile based on 3.x to 4.x based on the [official required changes](https://www.opendevstack.org/ods-documentation/opendevstack/4.x/update-guides/4x.html).

## Download

Target a specific version and platform for the binary. For version `v0.1.0` and linux, you would:

```cli
  curl -L https://github.com/opendevstack/ods-jenkinsfile-converter/releases/download/v0.1.0/converter -o converter
  sudo mv converter /usr/local/bin
```

For other platforms like Windows or MacOS choose the binary from the links [here](https://github.com/opendevstack/ods-jenkinsfile-converter/releases).

## Usage

Given a Jenkinsfile based on 3.x, to preview the changes that will result when migrating to ODS 4.x, use the `--dry-run` flag

```cli
  converter -filename="examples/machine-learning/Jenkinsfile" --dry-run
```

or to persist the changes in an output file

```cli
  converter -filename="examples/machine-learning/Jenkinsfile" -out "out/Jenkinsfile"
```

The converted Jenkinsfile with the ODS 4.x adopted changes will be under the [out](./out) directory:

```cli
  cat out/Jenkinsfile
```

## Building it locally

Generate the binaries under [bin](./bin)

```cli
  make all
```

You can run the binary for your platform. To preview the changes that will be done, use the `--dry-run` flag

```cli
  ./bin/converter -filename="examples/machine-learning/Jenkinsfile" --dry-run
```

## Run it with `Go`

```cli
  go run . -filename="examples/machine-learning/Jenkinsfile" -out "out/Jenkinsfile"
```

## Features

This tool will carry out the following changes in an ODS 3.x Jenkinsfile:

- [X] Update `@Library('ods-jenkins-shared-library@3.x') _` to `@Library('ods-jenkins-shared-library@4.x') _`
- [X] Point to agent images with the 4.x tag, e.g. change imageStreamTag: `ods/jenkins-agent-golang:3.x` to imageStreamTag: `ods/jenkins-agent-golang:4.x`.
- [X] `odsComponentStageImportOpenShiftImageOrElse` has been deprecated, and is now aliased to the new method, `odsComponentFindOpenShiftImageOrElse`.
- [X] `odsComponentStageRolloutOpenShiftDeployment` rolls out all deployment resources together now. If you had multiple DeploymentConfig resources previously, you had to target each one by specifying the config option resourceName. This is no longer possible - instead the stage iterates over all DeploymentConfig resources with the component label (app=${projectId}-${componentId}). Changes must be made to pipelines that have multiple deployments, such as components based on the ds-jupyter-notebook and ds-rshiny quickstarter.

  For instance, the following must be changed from:
  
    `odsComponentStageRolloutOpenShiftDeployment(context, [resourceName: "${context.componentId}-auth-proxy"])`
  
  to:

    `odsComponentStageRolloutOpenShiftDeployment(context)`
  
## Limitations
- Only Jenkinsfiles from 3.x are supported.
