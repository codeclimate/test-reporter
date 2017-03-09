package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cc-test-reporter",
	Short: "Report information about tests to Code Climate",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		v, err := cmd.Flags().GetBool("version")
		if err != nil {
			return err
		}
		if v {
			fmt.Printf("Code Climate Test Reporter %s (%s @ %s)\n", Version, BuildVersion, BuildTime)
		}
		return nil
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	RootCmd.Flags().BoolP("version", "v", false, "Show version information")
}
