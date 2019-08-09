package config

import "github.com/fatih/structs"

// ScriptFilter is an Alfred filter that runs a script
type ScriptFilter struct {
	Argument       keywordArgumentType `yaml:"argument" structs:"argumenttype"`
	Keyword        string              `yaml:"keyword" structs:"keyword"`
	RunningSubtext string              `yaml:"running-subtext" structs:"runningsubtext"`
	Title          string              `yaml:"title" structs:"title"`
	WithSpace      bool                `yaml:"with-space" structs:"withspace"`
	Script         ScriptConfig        `yaml:"script" structs:"-"`
}

func (s ScriptFilter) ToWorkflowConfig() map[string]interface{} {
	m := structs.Map(s)
	sMap := s.Script.ToWorkflowConfig()

	for k, v := range sMap {
		m[k] = v
	}

	return m
}
