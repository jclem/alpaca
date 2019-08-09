package config

import "github.com/fatih/structs"

// Clipboard is an object that copies to the clipboard
type Clipboard struct {
	Text string `yaml:"text" structs:"text"`
}

func (c Clipboard) ToWorkflowConfig() map[string]interface{} {
	return structs.Map(c)
}
