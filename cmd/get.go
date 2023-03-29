package cmd

import (
	"github.com/martinclaus1/zeus-client/pkg"
	"github.com/martinclaus1/zeus-client/pkg/loginPage"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"
)

var getCmd = &cobra.Command{
	Use:       "get status|overview",
	Short:     "Gets information from ZEUS® time tracking tool.",
	ValidArgs: []string{"status", "overview"},
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
		arg := args[0]
		defer pkg.Measure(time.Now(), "get "+arg)
		headless, _ := rootCmd.PersistentFlags().GetBool("headless")
		username, _ := rootCmd.PersistentFlags().GetString("username")
		password, _ := rootCmd.PersistentFlags().GetString("password")

		if "status" == arg {
			page := pkg.GetPage(&headless)
			currentStatus := loginPage.Instance(&page).Login(username, password).NavigateToMyZeusView().GetStatus()
			log.Infof("Current status: %s", currentStatus)
		}

		if "overview" == arg {
			page := pkg.GetPage(&headless)
			overview, _ := loginPage.Instance(&page).Login(username, password).NavigateToMyZeusView().GetOverview()
			overview.Print()
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
