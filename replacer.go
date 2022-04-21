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

// This attempts to change images that follows the pattern "ods/jenkins-agent*". The method considers single and double quoted
// values as well as prevents from panicking. Instead, the image should be left as is and it will advice to manually check it
// Note: this is considered a temporary solution to the problem faced with Jenkinsfiles
func replaceAgentImages(content string) string {

	reSingle := regexp.MustCompile(`'ods/jenkins-agent-(.*):.*'`)
	matches := reSingle.FindAllStringSubmatch(content, -1)

	if matches != nil {
		if string(matches[0][1]) == "nodejs10-angular" {
			return reSingle.ReplaceAllString(content, "'ods/jenkins-agent-nodejs12:4.x'")
		}

		return reSingle.ReplaceAllString(content, "'ods/jenkins-agent-$1:4.x'")

	} else {

		reDouble := regexp.MustCompile(`"ods/jenkins-agent-(.*):.*"`)
		matches := reDouble.FindAllStringSubmatch(content, -1)

		if matches != nil {
			if string(matches[0][1]) == "nodejs10-angular" {
				return reDouble.ReplaceAllString(content, "\"ods/jenkins-agent-nodejs12:4.x\"")
			}

			return reDouble.ReplaceAllString(content, "\"ods/jenkins-agent-$1:4.x\"")
		}
	}

	fmt.Printf("Warning: No canonical image found, such as this pattern :'ods/jenkins-agent-*'. \nPlease, consider if there is an image that should also be migrated to ODS version 4.\n")
	return content
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
