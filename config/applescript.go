package config

import "github.com/fatih/structs"

// AppleScript is an Alfred action that runs NSAppleScript
type AppleScript struct {
	Cache  bool         `yaml:"cache" structs:"cachescript"`
	Script ScriptConfig `yaml:"script" structs:"-"`
}

func (a AppleScript) ToWorkflowConfig() map[string]interface{} {
	m := structs.Map(a)
	m["applescript"] = a.Script.Content
	return m
}
