package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"
	"zeus-client/pkg"
	"zeus-client/pkg/loginPage"
)

var getCmd = &cobra.Command{
	Use:       "get status",
	Short:     "Gets information from ZEUS® time tracking tool.",
	ValidArgs: []string{"status"},
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.OnlyValidArgs(cmd, args); err != nil {
			return err
		}
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		defer pkg.Measure(time.Now(), "get status")
		headless, _ := rootCmd.PersistentFlags().GetBool("headless")
		username, _ := rootCmd.PersistentFlags().GetString("username")
		password, _ := rootCmd.PersistentFlags().GetString("password")

		if "status" == args[0] {
			page := pkg.GetPage(&headless)
			currentStatus := loginPage.Instance(&page).Login(username, password).NavigateToMyZeusView().GetStatus()
			log.Infof("Current status: %s", currentStatus)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
