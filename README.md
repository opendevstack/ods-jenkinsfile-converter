# ods-jenkinsfile-converter
[![build](https://github.com/opendevstack/ods-jenkinsfile-converter/actions/workflows/build.yml/badge.svg)](https://github.com/opendevstack/ods-jenkinsfile-converter/actions/workflows/build.yml)

## Introduction

This is a Golang CLI to convert the a Jenkinsfile based on 3.x to 4.x based on the [official required changes](https://www.opendevstack.org/ods-documentation/opendevstack/4.x/update-guides/4x.html).

## Example

Generate the binaries under [bin](./bin)

```cli
  make all
```

You can run the binary for your platform. To preview the changes that will be done, use the `--dry-run` flag

```cli
  ./bin/converter -filename="examples/machine-learning/Jenkinsfile" --dry-run"
```

or to persist the changes in an output file

```cli
  ./bin/converter -filename="examples/machine-learning/Jenkinsfile" -out "out/Jenkinsfile"
```

The converted Jenkinsfile with the ODS 4.x adopted changes is under the [out](./out) directory:

```cli
  cat out/Jenkinsfile
```

## Run it with `Go`

```cli
  go run . -filename="examples/machine-learning/Jenkinsfile" -out out/Jenkinsfile"
```

Considerations:

- It is assumed the Jenkins is based on ODS 3.x
- Update `@Library('ods-jenkins-shared-library@3.x') _` to `@Library('ods-jenkins-shared-library@4.x') _`
- Point to agent images with the 4.x tag, e.g. change imageStreamTag: `ods/jenkins-agent-golang:3.x` to imageStreamTag: `ods/jenkins-agent-golang:4.x`.
- `odsComponentStageImportOpenShiftImageOrElse` has been deprecated, and is now aliased to the new method, `odsComponentFindOpenShiftImageOrElse`.

- `odsComponentStageRolloutOpenShiftDeployment` rolls out all deployment resources together now. If you had multiple DeploymentConfig resources previously, you had to target each one by specifying the config option resourceName. This is no longer possible - instead the stage iterates over all DeploymentConfig resources with the component label (app=${projectId}-${componentId}). Changes must be made to pipelines that have multiple deployments, such as components based on the ds-jupyter-notebook and ds-rshiny quickstarter.

  For instance, the following must be changed from:
  
    `odsComponentStageRolloutOpenShiftDeployment(context, [resourceName: "${context.componentId}-auth-proxy"])`
  
  to:

    `odsComponentStageRolloutOpenShiftDeployment(context)`
