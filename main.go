package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	data := struct {
		DockerUrl   string `json:"docker_url"`
		Repository  string `json:"repository"`
		UpdatedTags string `json:"updated_tags"`
		Namespace   string `json:"namespace"`
		Name        string `json:"name"`
	}{}

	decoder := json.NewDecoder(os.Stdin)
	decoder.Decode(&data)

	fmt.Printf("data: %+v\n", data)
}
