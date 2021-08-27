package main

import "strings"

func valid(content string) bool {

	// Sanity check! Only Jenkinsfiles from 3.x are valid!
	return strings.Contains(content, "ods-jenkins-shared-library@3.x")
}
