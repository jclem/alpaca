package config

import (
	"io/ioutil"
	"os"

	"github.com/groob/plist"
	"gopkg.in/yaml.v2"
)

// Config is a parsed alpaca.json file.
type Config struct {
	Author              string                  `yaml:"author" plist:"createdby"`
	BundleID            string                  `yaml:"bundle-id" plist:"bundleid"`
	Connections         map[string][]Connection `yaml:"connections" plist:"connections"`
	Description         string                  `yaml:"description" plist:"description"`
	Name                string                  `yaml:"name" plist:"name"`
	Objects             map[string]Object       `yaml:"objects" plist:"-"`
	ObjectList          []Object                `yaml:"-" plist:"objects"`
	Readme              string                  `yaml:"readme" plist:"readme"`
	UIData              map[string]UIData       `yaml:"uidata" plist:"uidata"`
	URL                 string                  `yaml:"url" plist:"webaddress"`
	Variables           map[string]string       `yaml:"variables" plist:"variables"`
	VariablesDontExport []string                `yaml:"variablesdontexport" plist:"variablesdontexport"`
	Version             string                  `yaml:"version" plist:"version"`
}

// Connection is a line between two objects.
type Connection struct {
	To              string `yaml:"to" plist:"destinationuid"`
	Modifiers       int64  `yaml:"modifiers" plist:"modifiers"`
	ModifierSubtext string `yaml:"subtext" plist:"modifiersubtext"`
	VetoClose       bool   `yaml:"vito-close" plist:"vitoclose"` // NOTE: Yes, "vitoclose"
}

// Object is an object in an Alfred workflow.
type Object struct {
	Config  map[string]interface{} `yaml:"config" plist:"config"`
	Type    string                 `yaml:"type" plist:"type"`
	UID     string                 `yaml:"uid" plist:"uid"`
	Version int64                  `yaml:"version" plist:"version"`
}

// UIData represents the position of an object.
type UIData struct {
	X int64 `yaml:"x" plist:"xpos"`
	Y int64 `yaml:"y" plist:"ypos"`
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

// ToXML writes a config to XML.
func (c Config) ToXML() ([]byte, error) {
	// Construct ObjectList before we marshal to a plist.
	c.ObjectList = []Object{}
	for _, obj := range c.Objects {
		c.ObjectList = append(c.ObjectList, obj)
	}

	return plist.MarshalIndent(c, "\t")
}
