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
	}{}

	decoder := json.NewDecoder(os.Stdin)
	decoder.Decode(&data)

	script := `cd $(mktemp -d)
git clone https://github.com/boson-project/func/
cd func
go install ./cmd/func
cd ..
rm -fr func
`
	cmd := exec.Command("/bin/bash")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = strings.NewReader(script)

	err := cmd.Run()
	if err != nil {
		os.Exit(1)
	}

	fmt.Printf("data: %+v\n", data)
}


