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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			Replace(tt.args.inputFile, tt.args.convertedFile, false)

			converted, err := ioutil.ReadFile(tt.args.convertedFile)
			if err != nil {
				t.Errorf("failed to open file: %w", err)
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
