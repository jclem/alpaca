package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jclem/alpaca/project"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var out string

func init() {
	packCmd.Flags().StringVarP(&out, "out", "o", "", "Directory to output the packaged workflow to")
	rootCmd.AddCommand(&packCmd)
}

var packCmd = cobra.Command{
	Use:   "pack <dir>",
	Short: "Package the given Alpaca project into an Alfred workflow",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]

		projectPath, err := filepath.Abs(dir)
		if err != nil {
			log.Fatalf("Could not resolve path %s", dir)
		}

		outDir := out
		if outDir == "" {
			outDir, err = os.Getwd()
			if err != nil {
				err := errors.Wrap(err, "Error getting working directory")
				log.Fatal(err)
			}
		}

		project.Build(projectPath, outDir)
	},
}
