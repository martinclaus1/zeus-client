package cmd

import (
	"github.com/spf13/cobra"
	"zeus-client/pkg"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Sets up the username and password for the ZEUS® time tracking tool.",
	Long: `Sets up the username and password for the ZEUS® time tracking tool. 
The username and the password will be stored in an encrypted file. 
Therefore, a folder '.zeus-client' will be created in your home directory. 
Running the setup command will overwrite the existing credentials. 
The machine id is used as the encryption key.`,
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := rootCmd.PersistentFlags().GetString("username")
		password, _ := rootCmd.PersistentFlags().GetString("password")
		config := pkg.Config{Username: username, Password: password}
		config.Save()
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
