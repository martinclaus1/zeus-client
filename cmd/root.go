package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"net/mail"
	"os"
	"strings"
	"syscall"
	"zeus-client/pkg"
	"zeus-client/pkg/logging"
)

var debugMode bool
var username string
var password string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zeus-client",
	Short: "zeus-client is a CLI for the ZEUS time tracking tool",
	Long: `zeus-client is a CLI for the ZEUS time tracking tool. 
It allows you to toggle your presence state and get information about your current presence state, for example:
zeus-client toggle [flags]
zeus-client get status [flags]
zeus-client setup [flags]`,
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

		config := pkg.ReadConfig()

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
			bytepw, err := term.ReadPassword(int(syscall.Stdin))
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
	}
}
