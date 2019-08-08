package project

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/groob/plist"
	"github.com/jclem/alpaca/app/config"
	"github.com/jclem/alpaca/app/workflow"
	"github.com/pkg/errors"
)

// Build builds an Alpaca project into the given targetPath
func Build(projectDir string, targetDir string) error {
	cfg, err := readConfig(projectDir)
	if err != nil {
		return errors.Wrap(err, "Unable to read project config")
	}

	targetPath := filepath.Join(targetDir, fmt.Sprintf("%s.alfredworkflow", cfg.Name))

	workflowFile, err := os.Create(targetPath)
	if err != nil {
		return errors.Wrap(err, "Error creating workflow package file")
	}
	defer workflowFile.Close()

	archive := zip.NewWriter(workflowFile)
	defer archive.Close()

	writer, err := archive.Create("info.plist")
	if err != nil {
		return errors.Wrap(err, "Error creating plist file in archive")
	}

	info, err := workflow.NewFromConfig(projectDir, *cfg)
	if err != nil {
		return errors.Wrap(err, "Error creating worfklow from configuration")
	}
	plistBytes, err := plist.MarshalIndent(info, "\t")
	if err != nil {
		return errors.Wrap(err, "Error marshalling info plist")
	}

	writer.Write(plistBytes)

	if cfg.Icon != "" {
		src := filepath.Join(projectDir, cfg.Icon)
		ext := filepath.Ext(src)

		if ext != ".png" {
			return fmt.Errorf("Workflow icon must be a .png, got %q", ext)
		}

		dst := fmt.Sprintf("%s%s", "icon", ext)
		if err := copyFile(src, dst, archive); err != nil {
			return errors.Wrap(err, "Error copying file")
		}
	}

	for _, obj := range cfg.Objects {
		if obj.Icon == "" {
			continue
		}

		src := filepath.Join(projectDir, obj.Icon)
		ext := filepath.Ext(src)
		if ext != ".png" {
			return fmt.Errorf("Object icon must be a .png, got %q", ext)
		}

		dst := fmt.Sprintf("%s%s", obj.UID, ext)
		if err := copyFile(src, dst, archive); err != nil {
			return errors.Wrap(err, "Error copying file")
		}
	}

	if err := filepath.Walk(projectDir, func(filePath string, info os.FileInfo, err error) error {
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
		name := strings.TrimPrefix(filePath, projectDir+"/")
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
		return errors.Wrap(err, "Unable to create archive")
	}

	return nil
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
