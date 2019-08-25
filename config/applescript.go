package config

import (
	"github.com/fatih/structs"
	yaml "gopkg.in/yaml.v3"
)

// AppleScript is an Alfred action that runs NSAppleScript
type AppleScript struct {
	Cache   bool         `yaml:"cache" structs:"cachescript"`
	Content string       `yaml:"content" structs:"-"`
	Script  ScriptConfig `yaml:"script" structs:"-"`
}

func (a *AppleScript) UnmarshalYAML(node *yaml.Node) error {
	type alias AppleScript
	as := alias{Cache: true}
	if err := node.Decode(&as); err != nil {
		return err
	}

	*a = AppleScript(as)

	return nil
}

func (a AppleScript) ToWorkflowConfig() map[string]interface{} {
	m := structs.Map(a)
	m["applescript"] = a.Script.Content
	m["applescript"] = a.Content
	return m
}
