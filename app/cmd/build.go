package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jclem/alpaca/app/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build <dir>",
	Short: "Build the given directory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]
		path, err := filepath.Abs(dir)
		if err != nil {
			log.Fatalf("Could not resolve path %s", dir)
		}

		cfg, err := readConfig(path)
		if err != nil {
			log.Fatal(err)
		}

		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Could not get working directory")
		}

		targetPath := filepath.Join(wd, fmt.Sprintf("%s.alfredworkflow", cfg.Name))
		workflowFile, err := os.Create(targetPath)
		if err != nil {
			log.Fatal("Could not create workflow file")
		}
		defer workflowFile.Close()

		archive := zip.NewWriter(workflowFile)
		defer archive.Close()

		writer, err := archive.Create("info.plist")
		if err != nil {
			log.Fatal(err)
		}

		xml, err := cfg.ToXML()
		if err != nil {
			log.Fatal(err)
		}

		writer.Write(xml)

		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(path, "info.plist") {
				return nil
			}

			if err != nil {
				return err
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			writer, err := archive.CreateHeader(header)
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			if _, err := io.Copy(writer, file); err != nil {
				return err
			}

			return nil
		}); err != nil {
			log.Fatalf("Unable to create archive")
		}
	},
}

func readConfig(dir string) (*config.Config, error) {
	// Try .yaml first
	filePath := filepath.Join(dir, "alpaca.yaml")
	cfg, err := config.Read(filePath)

	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			// Then try .yml
			filePath = filepath.Join(dir, "alpaca.yml")
			cfg, err = config.Read(filePath)
			if err != nil {
				return nil, err
			}

			return cfg, nil
		}

		return nil, err
	}

	return cfg, nil
}
