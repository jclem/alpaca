package config

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

// Config is a parsed alpaca.json file.
type Config struct {
	Author      string
	BundleID    string            `yaml:"bundle-id"`
	Description string            `yaml:"description"`
	Icon        string            `yaml:"icon"`
	Name        string            `yaml:"name"`
	Objects     ObjectMap         `yaml:"objects"`
	Readme      string            `yaml:"readme"`
	URL         string            `yaml:"url"`
	Variables   map[string]string `yaml:"variables"`
	Version     string            `yaml:"version"`
}

// ObjectMap is a mapping of object names to objects
type ObjectMap map[string]Object

// Then is an object following another object.
type Then struct {
	Object string `yaml:"object"`
}

// UnmarshalYAML unmarshals an object.
func (o *ObjectMap) UnmarshalYAML(node *yaml.Node) error {
	var m map[string]Object
	if err := node.Decode(&m); err != nil {
		return err
	}

	*o = make(ObjectMap)

	for name, obj := range m {
		uid, err := uuid.NewRandom()
		if err != nil {
			return err
		}

		obj.UID = strings.ToUpper(uid.String())
		obj.Name = name
		(*o)[obj.Name] = obj
	}

	return nil
}

// Read parses an alpaca.json file.
func Read(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
