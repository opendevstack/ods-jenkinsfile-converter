package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

func Replace(filename, out string, dryRun bool) {

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	content := string(b)

	if !valid(content) {
		log.Fatalf("The Jenkinsfile %s is not based on 3.x and thus cannot be converted to 4.x.", filename)
	}

	content = replaceLibrary(content)
	content = replaceAgentImages(content)
	content = replaceComponentStageImport(content)
	content = replaceComponentStageRollout(content)

	if dryRun {
		fmt.Println(content)
	} else {
		ioutil.WriteFile(out, []byte(content), 0666)
		log.Printf("Generated file %s", out)
	}
}

func replaceLibrary(content string) string {
	return strings.Replace(content, "@Library('ods-jenkins-shared-library@3.x') _", "@Library('ods-jenkins-shared-library@4.x') _", -1)
}

func replaceAgentImages(content string) string {
	re := regexp.MustCompile(`'ods/jenkins-agent-(.*):.*'`)

	matches := re.FindAllStringSubmatch(content, -1)
	fmt.Println(matches[0][1])
	if string(matches[0][1]) == "nodejs10-angular" {
		return re.ReplaceAllString(content, "'ods/jenkins-agent-nodejs12:4.x'")
	}

	return re.ReplaceAllString(content, "'ods/jenkins-agent-$1:4.x'")
}

func replaceComponentStageImport(content string) string {
	return strings.Replace(content, "odsComponentStageImportOpenShiftImageOrElse", "odsComponentFindOpenShiftImageOrElse", -1)
}

func replaceComponentStageRollout(content string) string {
	content, match := replaceComponentStageRolloutMultiLine(content)
	if !match {
		content, _ = replaceComponentStageRolloutSingleLine(content)
	}
	return content
}

func replaceComponentStageRolloutMultiLine(content string) (string, bool) {
	re := regexp.MustCompile(`(?ms)odsComponentStageRolloutOpenShiftDeployment\((context, \[\n?.*)]\)$`)
	return re.ReplaceAllString(content, "odsComponentStageRolloutOpenShiftDeployment(context)"), re.Match([]byte(content))
}

func replaceComponentStageRolloutSingleLine(content string) (string, bool) {
	re := regexp.MustCompile(`odsComponentStageRolloutOpenShiftDeployment(.*)`)
	return re.ReplaceAllString(content, "odsComponentStageRolloutOpenShiftDeployment(context)"), re.Match([]byte(content))
}
