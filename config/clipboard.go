package config

import (
	"github.com/fatih/structs"
	yaml "gopkg.in/yaml.v3"
)

// Clipboard is an object that copies to the clipboard
type Clipboard struct {
	Text string `yaml:"text" structs:"clipboardtext"`
}

func (c *Clipboard) UnmarshalYAML(node *yaml.Node) error {
	type alias Clipboard
	as := alias{Text: "{query}"}
	if err := node.Decode(&as); err != nil {
		return err
	}

	*c = Clipboard(as)

	return nil
}

func (c Clipboard) ToWorkflowConfig() map[string]interface{} {
	return structs.Map(c)
}
