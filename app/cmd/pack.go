package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/groob/plist"
	"github.com/jclem/alpaca/app/config"
	"github.com/jclem/alpaca/app/workflow"
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
		path, err := filepath.Abs(dir)
		if err != nil {
			log.Fatalf("Could not resolve path %s", dir)
		}

		cfg, err := readConfig(path)
		if err != nil {
			log.Fatal(err)
		}

		outDir := out

		if out == "" {
			outDir, err = os.Getwd()
			if err != nil {
				log.Fatalf("Could not get working directory")
			}
		}

		targetPath := filepath.Join(outDir, fmt.Sprintf("%s.alfredworkflow", cfg.Name))
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

		info, err := workflow.NewFromConfig(path, *cfg)
		if err != nil {
			log.Fatal(err)
		}
		plistBytes, err := plist.MarshalIndent(info, "\t")
		if err != nil {
			log.Fatal(err)
		}

		writer.Write(plistBytes)

		if cfg.Icon != "" {
			src := filepath.Join(path, cfg.Icon)
			ext := filepath.Ext(src)

			if ext != ".png" {
				log.Fatalf("Workflow icon must be a .png, got %q", ext)
			}

			dst := fmt.Sprintf("%s%s", "icon", ext)
			if err := copyFile(src, dst, archive); err != nil {
				log.Fatal(err)
			}
		}

		for _, obj := range cfg.Objects {
			if obj.Icon == "" {
				continue
			}

			src := filepath.Join(path, obj.Icon)
			ext := filepath.Ext(src)
			if ext != ".png" {
				log.Fatalf("Object icon must be a .png, got %q", ext)
			}

			dst := fmt.Sprintf("%s%s", obj.UID, ext)
			if err := copyFile(src, dst, archive); err != nil {
				log.Fatal(err)
			}
		}

		if err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			if filePath == targetPath {
				return nil
			}

			if strings.HasSuffix(filePath, "info.plist") {
				return nil
			}

			if err != nil {
				return err
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}
			name := strings.TrimPrefix(filePath, path+"/")
			header.Name = name

			writer, err := archive.CreateHeader(header)
			if err != nil {
				return err
			}

			file, err := os.Open(filePath)
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

func copyFile(from string, to string, archive *zip.Writer) error {
	info, err := os.Stat(from)
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	header.Name = to

	writer, err := archive.CreateHeader(header)
	if err != nil {
		return err
	}

	file, err := os.Open(from)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(writer, file); err != nil {
		return err
	}

	return nil
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
