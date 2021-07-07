package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	data := struct {
		DockerUrl   string   `json:"docker_url"`
		Repository  string   `json:"repository"`
		UpdatedTags []string `json:"updated_tags"`
		Namespace   string   `json:"namespace"`
		Name        string   `json:"name"`
		Runtimes    []string `json:"runtimes"`
	}{}

	decoder := json.NewDecoder(os.Stdin)
	decoder.Decode(&data)

	if len(data.UpdatedTags) < 1 {
		fmt.Fprintf(os.Stderr, "there are no updated tags")
		os.Exit(1)
	}

	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get working directory: %q", err.Error())
		os.Exit(0)
	}
	path := os.Getenv("PATH")
	os.Setenv("PATH", fmt.Sprintf("%s:%s", workingDir, path))

	os.Setenv("FUNC_REGISTRY", "example.com/jdoe")

	err = installFunc(workingDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to install func: %q\n", err.Error())
		os.Exit(1)
	}

	builderImg := data.DockerUrl + ":" + data.UpdatedTags[0]

	failed := false

	for _, funcBinary := range []string {"func_stable", "func_latest"} {
		for _, runtime := range data.Runtimes {
			for _, template := range []string {"http", "events"} {
				err = tryBuild(funcBinary, runtime, template, builderImg)
				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to build the function (func binary: %q, template: %q): %q\n", funcBinary, template, err.Error())
					failed = true
				}
			}
		}
	}
	if failed {
		os.Exit(1)
	}
}

// installs two func binaries to the target dir
// func_stable -- latest tagged version
// func_latest -- latest main
func installFunc(trg string) error {
	script := fmt.Sprintf(`set -ex
go get github.com/markbates/pkger/cmd/pkger
cd $(mktemp -d)
git clone https://github.com/boson-project/func
cd func
make
cp func %[1]s/func_latest
make clean
git fetch --tags
latestTag=$(git describe --tags $(git rev-list --tags --max-count=1))
git checkout $latestTag
make
cp func %[1]s/func_stable
cd ..
rm -fr func`, trg)
	return runBash(script)
}

// creates and tries to build a function
func tryBuild(funcBinary string, runtime string, template string, builderImg string) error {
	script := fmt.Sprintf(`set -ex
cd $(mktemp -d)
%[1]s create fn%[2]s%[3]s --runtime %[2]s --template %[3]s
cd fn%[2]s%[3]s
%[1]s build --builder %[4]s`, funcBinary, runtime, template, builderImg)
	return runBash(script)
}

func runBash(in string) error {
	cmd := exec.Command("/bin/bash")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = strings.NewReader(in)
	return cmd.Run()
}


