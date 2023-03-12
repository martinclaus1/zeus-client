package cmd

import (
	"fmt"
	buildConfig "github.com/martinclaus1/zeus-client/config"
	"github.com/martinclaus1/zeus-client/pkg"
	"github.com/martinclaus1/zeus-client/pkg/logging"
	"github.com/martinclaus1/zeus-client/pkg/user"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"golang.org/x/term"
	"net/mail"
	"os"
	"strings"
	"syscall"
)

var debugMode bool
var username string
var password string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zeus-client",
	Short: "zeus-client is a CLI for the ZEUS® time tracking tool.",
	Long: `zeus-client is a CLI for the ZEUS® time tracking tool. It allows you to track your time and to show information using the command line. 
Playwright is used to automate the browser interaction. Credentials can be saved using the 'setup' command.`,
	DisableAutoGenTag: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "Enable debug mode")
	rootCmd.PersistentFlags().Bool("headless", true, "Runs the playwright script in headless mode")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "Username for the ZEUS time tracking tool")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Password for the ZEUS time tracking tool. If not provided, the script will prompt for the password")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		logging.SetupLogging(debugMode)
		pkg.InstallPlaywright()

		config := user.ReadConfig()

		// try username from config if username is not provided
		if username == "" {
			username = config.Username
		}

		if username == "" {
			logrus.Fatalln("Username must be provided.")
		}

		if _, err := mail.ParseAddress(username); err != nil {
			logrus.Fatalln("Username is not valid ", err)
		}

		// try password from config if password is not provided
		if password == "" {
			password = config.Password
		}

		if password == "" {
			fmt.Print("Password: ")
			bytepw, err := term.ReadPassword(syscall.Stdin)
			if err != nil {
				logrus.Fatalln("Could not read password")
			}
			p := strings.TrimSpace(string(bytepw))
			fmt.Println()

			if p == "" {
				logrus.Fatalln("Password must be provided.")
			}

			if err = rootCmd.PersistentFlags().Set("password", p); err != nil {
				logrus.WithField("error", err).Fatalln("Could not set password.")
			}
		}

		logrus.Debugf("Version: %s, build time: %s", buildConfig.Version, buildConfig.BuildTime)
	}
}

func GenerateDocumentation(dir string) {
	err := doc.GenMarkdownTree(rootCmd, dir)
	if err != nil {
		logrus.Fatalln("could not generate documentation: %w", err)
	}
}
