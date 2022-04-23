package main

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestReplace(t *testing.T) {

	type args struct {
		inputFile     string
		convertedFile string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Stage Rollout Single Line",
			args: args{
				inputFile:     "internal/test/stage-rollout-single-line/Jenkinsfile",
				convertedFile: "internal/test/stage-rollout-single-line/out/Jenkinsfile",
			},
			want: readFile("internal/test/stage-rollout-single-line/golden/Jenkinsfile"),
		},
		{
			name: "Stage Rollout Multi Line",
			args: args{
				inputFile:     "internal/test/stage-rollout-multi-line/Jenkinsfile",
				convertedFile: "internal/test/stage-rollout-multi-line/out/Jenkinsfile",
			},
			want: readFile("internal/test/stage-rollout-multi-line/golden/Jenkinsfile"),
		},
		{
			name: "NodeJS agent name change",
			args: args{
				inputFile:     "internal/test/nodejs-agent-name/Jenkinsfile",
				convertedFile: "internal/test/nodejs-agent-name/out/Jenkinsfile",
			},
			want: readFile("internal/test/nodejs-agent-name/golden/Jenkinsfile"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			Replace(tt.args.inputFile, tt.args.convertedFile, false)

			converted, err := ioutil.ReadFile(tt.args.convertedFile)
			if err != nil {
				t.Errorf("failed to open file: %v", err)
			}

			if got := string(converted); got != string(tt.want) {
				t.Errorf("Replace() =\n%v, want\n%v", got, string(tt.want))
			}
		})
	}

}

func readFile(goldenFile string) string {
	want, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		panic(err)
	}

	return string(want)
}

// The focus has been: matching patterns, quotes, duplicity of images, combination of them
func TestTableReplaceAgentImages(t *testing.T) {

	testsDouble := []struct {
		input    string
		expected string
	}{
		{input: "ods/jenkins-agent-base:3.0.0", expected: "\"ods/jenkins-agent-base:4.x\""},
		{input: "${dockerRegistry}/ods/jenkins-agent-nodejs10-angular:3.0.0", expected: "\"${dockerRegistry}/ods/jenkins-agent-nodejs10-angular:3.0.0\""},
		{input: "ods/jenkins-agent-nodejs10-angular:3.0.0", expected: "\"ods/jenkins-agent-nodejs12:4.x\""},
		{input: "odsalpha/jenkins-agent-nodejs10-angular:3.x", expected: "\"odsalpha/jenkins-agent-nodejs10-angular:3.x\""},
		{input: "odsalpha/jenkins-agent-base:3.x", expected: "\"odsalpha/jenkins-agent-base:3.x\""},
		{input: "", expected: "\"\""},
		{input: "${dockerRegistry}/edpp-cd/jenkins-agent-maven-chrome:latest", expected: "\"${dockerRegistry}/edpp-cd/jenkins-agent-maven-chrome:latest\""},
	}

	fmt.Println()
	t.Logf("#####################################")
	t.Logf("## Test suite double quoted values ##")
	t.Logf("#####################################\n")
	fmt.Println()
	for _, test := range testsDouble {
		doubleQuoted := fmt.Sprintf("%q", test.input)
		if output := replaceAgentImages(doubleQuoted); output != test.expected {
			t.Errorf("FAILED; input: %v, expected: %v, received: %v.", doubleQuoted, test.expected, output)
		} else {
			t.Logf("PASSED; input: %v, expected: %v, received: %v", doubleQuoted, test.expected, output)
		}
	}

	testsSingle := []struct {
		input    string
		expected string
	}{
		{input: "alpha/jenkins-agent-nodejs10-angular:3.x", expected: "'alpha/jenkins-agent-nodejs10-angular:3.x'"},
		{input: "ods/jenkins-agent-maven:3.0.0", expected: "'ods/jenkins-agent-maven:4.x'"},
		{input: "${dockerRegistry}/ods/jenkins-agent-nodejs10-angular:3.0.0", expected: "'${dockerRegistry}/ods/jenkins-agent-nodejs10-angular:3.0.0'"},
		{input: "ods/jenkins-agent-nodejs10-angular:3.0.0", expected: "'ods/jenkins-agent-nodejs12:4.x'"},
		{input: "", expected: "''"},
		{input: "odsalpha/jenkins-agent-terraform:3.x", expected: "'odsalpha/jenkins-agent-terraform:3.x'"},
		{input: "${dockerRegistry}/edpp-cd/jenkins-agent-maven-chrome:latest", expected: "'${dockerRegistry}/edpp-cd/jenkins-agent-maven-chrome:latest'"},
	}

	fmt.Println()
	fmt.Println()
	t.Logf("#####################################")
	t.Logf("## Test suite single quoted values ##")
	t.Logf("#####################################\n")
	fmt.Println()
	for _, test := range testsSingle {
		singleQuoted := fmt.Sprintf("'%v'", test.input)
		if output := replaceAgentImages(singleQuoted); output != test.expected {
			t.Errorf("FAILED; input: %v, expected: %v, received: %v.", singleQuoted, test.expected, output)
		} else {
			t.Logf("PASSED; input: %v, expected: %v, received: %v", singleQuoted, test.expected, output)
		}
	}

	testsMultiLines := []struct {
		input    string
		expected string
	}{
		{input: " 'alpha/jenkins-agent-nodejs10-angular:3.x'\n 'alpha/jenkins-agent-nodejs10-angular:3.x'\n", expected: " 'alpha/jenkins-agent-nodejs10-angular:3.x'\n 'alpha/jenkins-agent-nodejs10-angular:3.x'\n"},
		{input: " 'ods/jenkins-agent-maven:3.0.0'\n 'ods/jenkins-agent-maven:3.0.0'\n", expected: " 'ods/jenkins-agent-maven:4.x'\n 'ods/jenkins-agent-maven:4.x'\n"},
		{input: " 'ods/jenkins-agent-base:3.0.0'\n 'ods/jenkins-agent-base:4.x'\n", expected: " 'ods/jenkins-agent-base:4.x'\n 'ods/jenkins-agent-base:4.x'\n"},
		{input: " 'ods/jenkins-agent-base:4.x'\n 'ods/jenkins-agent-base:3.0.0'\n", expected: " 'ods/jenkins-agent-base:4.x'\n 'ods/jenkins-agent-base:4.x'\n"},
		{input: " '${dockerRegistry}/ods/jenkins-agent-nodejs10-angular:3.0.0'\n '${dockerRegistry}/ods/jenkins-agent-nodejs10-angular:3.0.0'\n", expected: " '${dockerRegistry}/ods/jenkins-agent-nodejs10-angular:3.0.0'\n '${dockerRegistry}/ods/jenkins-agent-nodejs10-angular:3.0.0'\n"},
		{input: " '${dockerRegistry}/ods/jenkins-agent-nodejs10-angular:3.0.0'\n '${dockerRegistry}/ods/jenkins-agent-nodejs12-angular:4.x'\n", expected: " '${dockerRegistry}/ods/jenkins-agent-nodejs10-angular:3.0.0'\n '${dockerRegistry}/ods/jenkins-agent-nodejs12-angular:4.x'\n"},
		{input: " 'ods/jenkins-agent-nodejs10-angular:3.0.0'\n 'ods/jenkins-agent-nodejs12:4.x'\n", expected: " 'ods/jenkins-agent-nodejs12:4.x'\n 'ods/jenkins-agent-nodejs12:4.x'\n"},
		{input: " 'ods/jenkins-agent-nodejs12:4.x'\n 'ods/jenkins-agent-nodejs10-angular:3.0.0'\n", expected: " 'ods/jenkins-agent-nodejs12:4.x'\n 'ods/jenkins-agent-nodejs12:4.x'\n"},
		{input: " '${dockerRegistry}/edpp-cd/jenkins-agent-maven-chrome:latest'\n '${dockerRegistry}/edpp-cd/jenkins-agent-maven-chrome:latest'\n", expected: " '${dockerRegistry}/edpp-cd/jenkins-agent-maven-chrome:latest'\n '${dockerRegistry}/edpp-cd/jenkins-agent-maven-chrome:latest'\n"},
	}

	fmt.Println()
	fmt.Println()
	t.Logf("######################################")
	t.Logf("## Test suite multiple lines values ##")
	t.Logf("######################################\n")
	fmt.Println()
	for _, test := range testsMultiLines {
		multiLine := fmt.Sprintf("%v", test.input)
		fmt.Println()
		if output := replaceAgentImages(multiLine); output != test.expected {
			t.Errorf("FAILED; input: %v, expected: %v, received: %v.", multiLine, test.expected, output)
		} else {
			t.Logf("PASSED; input: %v, expected: %v, received: %v", multiLine, test.expected, output)
		}
	}
}
