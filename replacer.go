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

// replaceAgentImages is a method that changes images which follows the pattern "ods/jenkins-agent*".
//It considers single, double quotes, and multiple images in the same Jenkinsfile without panicking.
func replaceAgentImages(content string) string {

	re := regexp.MustCompile(`(['"]{1}ods/jenkins-agent-)(.*):(.*)(['"]{1})`)
	lines := strings.Split(content, "\n")

	for idx, line := range lines {

		if match := re.FindStringSubmatch(line); match != nil {
			lines[idx] = func(s1 string, s2 []string) string {
				if string(match[2]) == "nodejs10-angular" {
					return re.ReplaceAllString(line, fmt.Sprintf("%s%s%s", s2[1], "nodejs12:4.x", s2[4]))
				} else {
					return re.ReplaceAllString(line, fmt.Sprintf("%s%s%s%s", s2[1], s2[2], ":4.x", s2[4]))
				}
			}(line, match)
		}
	}

	return strings.Join(lines, "\n")
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
