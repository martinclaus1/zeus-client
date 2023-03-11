package cmd

import (
	"github.com/martinclaus1/zeus-client/pkg"
	"github.com/martinclaus1/zeus-client/pkg/loginPage"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"
)

var toggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "Toggles the presence state in ZEUS® time tracking tool.",
	Long: `Toggles the presence state in ZEUS® time tracking tool.
If the current state is "present", it will be set to "absent" and vice versa. 
In the end, the current state will be printed to the console.`,
	Run: func(cmd *cobra.Command, args []string) {
		defer pkg.Measure(time.Now(), "toggle")
		headless, _ := rootCmd.PersistentFlags().GetBool("headless")
		username, _ := rootCmd.PersistentFlags().GetString("username")
		password, _ := rootCmd.PersistentFlags().GetString("password")

		page := pkg.GetPage(&headless)
		myZeusView := loginPage.Instance(&page).Login(username, password).NavigateToMyZeusView()
		log.Infof("Current status: %s", myZeusView.GetStatus())
		myZeusView.ToggleStatus()
	},
}

func init() {
	rootCmd.AddCommand(toggleCmd)
}
