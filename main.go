package main

import (
	"flag"
	"fmt"
)

func main() {

	filename := flag.String("filename", "Jenkinsfile", "Name of the Jenkinsfile")
	out := flag.String("out", "Jenkinsfile.4x", "Name of the Jenkinsfile with the adopted modifications for 4.x")
	dryRun := flag.Bool("dry-run", false, "to preview the modifications that would be done in a new Jenkinsfile, without really creating it")

	flag.Parse()

	fmt.Println("filename: ", *filename)
	fmt.Println("out: ", *out)
	fmt.Println("dry-run:", *dryRun)

	Replace(*filename, *out, *dryRun)
}
