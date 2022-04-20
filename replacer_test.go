package main

import (
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
		{
			name: "Stage Rollout Single Line",
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

// We expect to fail here almost in all cases, except for ods/jenkins-agent-* images
// It has also the problem of the single quotes in the pattern
func TestTableReplaceAgentImages(t *testing.T) {

	tests := []struct {
		input    string
		expected string
	}{
		{input: "ods/jenkins-agent-nodejs10-angular:3.0.0", expected: "'ods/jenkins-agent-nodejs12-angular:4.x'"},
		{input: "ods/jenkins-agent-maven:3.0.0", expected: "'ods/jenkins-agent-maven:4.x'"},
		{input: "ods/jenkins-agent-base:3.0.0", expected: "'ods/jenkins-agent-base:4.x'"},
		{input: "", expected: ""},
		{input: "alpha/jenkins-agent-nodejs10-angular:3.x", expected: "alpha/jenkins-agent-nodejs10-angular:3.x"},
		{input: "${dockerRegistry}/ods/jenkins-agent-nodejs10-angular:3.0.0", expected: "${dockerRegistry}/ods/jenkins-agent-nodejs10-angular:3.0.0"},
		{input: "odsalpha/jenkins-agent-nodejs10-angular:3.x", expected: "odsalpha/jenkins-agent-nodejs10-angular:3.x"},
		{input: "odsalpha/jenkins-agent-base:3.x", expected: "odsalpha/jenkins-agent-base:3.x"},
		{input: "odsalpha/jenkins-agent-terraform:3.x", expected: "odsalpha/jenkins-agent-terraform:3.x"},
		{input: "${dockerRegistry}/edpp-cd/jenkins-agent-maven-chrome:latest", expected: "${dockerRegistry}/edpp-cd/jenkins-agent-maven-chrome:latest"},
	}

	for _, test := range tests {
		if output := replaceAgentImages(test.input); output != test.expected {
			t.Errorf("TestTableReplaceAgentImages FAILED; input: %q, expected %q, received: %q", test.input, test.expected, output)
		} else {
			t.Logf("TestTableReplaceAgentImages PASSED; input: %q, expected %q, received: %q", test.input, test.expected, output)
		}
	}
}

// At the moment we expect three changes to take place in this test suite, one blank without panic (relman)
// and the rest to be left the same because they are concrete images
func TestTableReplaceAgentImagesFuzzy(t *testing.T) {

	tests := []struct {
		input    string
		expected string
	}{
		{input: "alpha/jenkins-agent-nodejs10-angular:3.x", expected: "alpha/jenkins-agent-nodejs10-angular:3.x"},
		{input: "ods/jenkins-agent-maven:3.0.0", expected: "ods/jenkins-agent-maven:4.x"},
		{input: "ods/jenkins-agent-base:3.0.0", expected: "ods/jenkins-agent-base:4.x"},
		{input: "${dockerRegistry}/ods/jenkins-agent-nodejs10-angular:3.0.0", expected: "${dockerRegistry}/ods/jenkins-agent-nodejs10-angular:3.0.0"},
		{input: "ods/jenkins-agent-nodejs10-angular:3.0.0", expected: "'ods/jenkins-agent-nodejs12:4.x'"},
		{input: "odsalpha/jenkins-agent-nodejs10-angular:3.x", expected: "odsalpha/jenkins-agent-nodejs10-angular:3.x"},
		{input: "odsalpha/jenkins-agent-base:3.x", expected: "odsalpha/jenkins-agent-base:3.x"},
		{input: "", expected: ""},
		{input: "odsalpha/jenkins-agent-terraform:3.x", expected: "odsalpha/jenkins-agent-terraform:3.x"},
		{input: "${dockerRegistry}/edpp-cd/jenkins-agent-maven-chrome:latest", expected: "${dockerRegistry}/edpp-cd/jenkins-agent-maven-chrome:latest"},
	}

	for _, test := range tests {
		if output := replaceAgentImagesFuzzy(test.input); output != test.expected {
			t.Errorf("TestTableReplaceAgentImages FAILED; input: %q, expected %q, received: %q.", test.input, test.expected, output)
		} else {
			t.Logf("TestTableReplaceAgentImages PASSED; input: %q, expected %q, received: %q", test.input, test.expected, output)
		}
	}

}
