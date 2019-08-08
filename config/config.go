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
	Author      string            `yaml:"author"`
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

// Object is an object in an Alfred workflow.
type Object struct {
	Argument  string  `yaml:"argument"`
	Icon      string  `yaml:"icon"`
	Keyword   string  `yaml:"keyword"`
	Name      string  `yaml:"-"`
	Script    *Script `yaml:"script"`
	Then      []Then  `yaml:"then"`
	Type      string  `yaml:"type"`
	UID       string  `yaml:"uid"`
	Version   int64   `yaml:"version"`
	WithSpace bool    `yaml:"with-space"`
}

// ObjectMap is a mapping of object names to objects
type ObjectMap map[string]Object

// Script is a runnable script in a workflow.
type Script struct {
	Content string `yaml:"content"`
	Inline  bool   `yaml:"inline"`
	Path    string `yaml:"path"`
	Type    string `yaml:"type"`
}

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
