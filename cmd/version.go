package cmd

import (
	"fmt"
	buildConfig "zeus-client/config"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of the ZEUS® time tracking tool.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s, build time: %s\n", buildConfig.Version, buildConfig.BuildTime)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
