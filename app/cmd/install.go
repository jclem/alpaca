package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(&installCmd)
}

var installCmd = cobra.Command{}
