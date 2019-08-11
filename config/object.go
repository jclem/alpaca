package config

import (
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

// ObjectType denotes the general behavior of an Alfred object
type ObjectType string

var (
	AppleScriptType  ObjectType = "applescript"
	ClipboardType    ObjectType = "clipboard"
	KeywordType      ObjectType = "keyword"
	OpenURLType      ObjectType = "open-url"
	ScriptType       ObjectType = "script"
	ScriptFilterType ObjectType = "script-filter"
	UnknownType      ObjectType = "unknown"
)

// Object is an object in an Alfred workflow
type Object struct {
	Name    string       `yaml:"-" structs:"-"`
	Icon    string       `yaml:"icon" structs:"-"`
	Type    ObjectType   `yaml:"type" structs:"-"`
	UID     string       `yaml:"-" structs:"uid"`
	Then    []Then       `yaml:"then" structs:"-"`
	Version int64        `yaml:"version" structs:"version"`
	Config  ObjectConfig `yaml:"config" structs:"-"`
}

func (o *Object) UnmarshalYAML(node *yaml.Node) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	o.UID = strings.ToUpper(uuid.String())

	var proxy struct {
		Icon    string
		Type    ObjectType
		Then    []Then
		Version int64
		Config  map[string]interface{}
	}

	if err := node.Decode(&proxy); err != nil {
		return err
	}

	o.Icon = proxy.Icon
	o.Type = proxy.Type
	o.Then = proxy.Then
	o.Version = proxy.Version

	rawConfig, err := yaml.Marshal(proxy.Config)
	if err != nil {
		return err
	}

	switch o.Type {
	case AppleScriptType:
		var cfg AppleScript
		if err := yaml.Unmarshal(rawConfig, &cfg); err != nil {
			return err
		}
		o.Config = cfg
	case ClipboardType:
		var cfg Clipboard
		if err := yaml.Unmarshal(rawConfig, &cfg); err != nil {
			return err
		}
		o.Config = cfg
	case KeywordType:
		var cfg Keyword
		if err := yaml.Unmarshal(rawConfig, &cfg); err != nil {
			return err
		}
		o.Config = cfg
	case OpenURLType:
		var cfg OpenURL
		if err := yaml.Unmarshal(rawConfig, &cfg); err != nil {
			return err
		}
		o.Config = cfg
	case ScriptType:
		var cfg Script
		if err := yaml.Unmarshal(rawConfig, &cfg); err != nil {
			return err
		}
		o.Config = cfg
	case ScriptFilterType:
		var cfg ScriptFilter
		if err := yaml.Unmarshal(rawConfig, &cfg); err != nil {
			return err
		}
		o.Config = cfg
	default:
		fmt.Println(o.Type)
	}

	return nil
}

var objectType = map[ObjectType]string{
	"applescript":   "alfred.workflow.action.applescript",
	"clipboard":     "alfred.workflow.output.clipboard",
	"keyword":       "alfred.workflow.input.keyword",
	"open-url":      "alfred.workflow.action.openurl",
	"script":        "alfred.workflow.action.script",
	"script-filter": "alfred.workflow.input.scriptfilter",
}

func (o Object) ToWorkflowConfig() map[string]interface{} {
	m := structs.Map(o)
	m["type"] = objectType[o.Type]
	m["config"] = o.Config.ToWorkflowConfig()
	return m
}

// ObjectConfig is a general configuration for an object
type ObjectConfig interface {
	ToWorkflowConfig() map[string]interface{}
}

var scriptType = map[string]int64{
	"bash":         0,
	"php":          1,
	"ruby":         2,
	"python":       3,
	"perl":         4,
	"zsh":          5,
	"osascript-as": 6,
	"osascript-js": 7,
	"external":     8,
}

// ScriptConfig is a runnable script in a workflow.
type ScriptConfig struct {
	Content string `yaml:"content" structs:"script"`
	Path    string `yaml:"path" structs:"scriptfile"`
	Type    string `yaml:"type" structs:"-"`
}

func (s ScriptConfig) ToWorkflowConfig() map[string]interface{} {
	m := structs.Map(s)

	if s.Type == "" {
		m["type"] = scriptType["external"]
	} else {
		m["type"] = scriptType[s.Type]
	}

	return m
}
