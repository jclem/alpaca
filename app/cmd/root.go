package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:   "alpaca",
	Short: "Alpaca is a packaging utility for Alfred workflows",
	Long: `A package utility for Alfred workflows built by @jclem in Go

Documentation at https://github.com/jclem/alpaca`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
