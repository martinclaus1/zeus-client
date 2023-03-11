package main

import (
	"github.com/martinclaus1/zeus-client/cmd"
	"os"
	"strings"
)

func main() {
	if len(os.Args) == 2 && strings.Contains(os.Args[0], "exe") && os.Args[1] == "generate-docs" {
		cmd.GenerateDocumentation("documentation")
	} else {
		cmd.Execute()
	}
}
