package main

import (
	"os"
	"strings"
	"zeus-client/cmd"
)

func main() {
	if len(os.Args) == 2 && strings.Contains(os.Args[0], "exe") && os.Args[1] == "generate-docs" {
		cmd.GenerateDocumentation("documentation")
	} else {
		cmd.Execute()
	}
}
