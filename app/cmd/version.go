package cmd

import (
	"fmt"

	"github.com/jclem/alpaca/app/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Alpaca",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Version)
	},
}
