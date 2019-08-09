package config

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

// ObjectType denotes the general behavior of an Alfred object
type ObjectType string

var (
	ClipboardType    ObjectType = "clipboard"
	KeywordType      ObjectType = "keyword"
	OpenURLType      ObjectType = "open-url"
	ScriptType       ObjectType = "script"
	ScriptFilterType ObjectType = "script-filter"
	UnknownType      ObjectType = "unknown"
)

// Object is an object in an Alfred workflow
type Object struct {
	Name    string       `yaml:"name" structs:"-"`
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
	o.UID = uuid.String()

	var proxy struct {
		Name    string
		Icon    string
		Type    ObjectType
		Then    []Then
		Version int64
		Config  map[string]interface{}
	}

	if err := node.Decode(&proxy); err != nil {
		return err
	}

	o.Name = proxy.Name
	o.Icon = proxy.Icon
	o.Type = proxy.Type
	o.Then = proxy.Then
	o.Version = proxy.Version

	rawConfig, err := yaml.Marshal(proxy.Config)
	if err != nil {
		return err
	}

	switch o.Type {
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

type keywordArgumentType string

var (
	keywordArgumentRequired keywordArgumentType = "required"
	keywordArgumentOptional                     = "optional"
	keywordArgumentNone                         = "none"
)

// Keyword is an object triggered by a keyword
type Keyword struct {
	Keyword   string              `yaml:"keyword" structs:"keyword"`
	WithSpace bool                `yaml:"with-space" structs:"withspace"`
	Argument  keywordArgumentType `yaml:"argument" structs:"argumenttype"`
}

var argumentType = map[keywordArgumentType]int64{
	"required": 0,
	"optional": 1,
	"none":     2,
}

func (k Keyword) ToWorkflowConfig() map[string]interface{} {
	m := structs.Map(k)
	m["argumenttype"] = argumentType[k.Argument]
	return m
}

// Script is an Alfred action that runs a script
type Script struct {
	Script ScriptConfig
}

func (s Script) ToWorkflowConfig() map[string]interface{} {
	return s.Script.ToWorkflowConfig()
}

// ScriptFilter is an Alfred filter that runs a script
type ScriptFilter struct {
	Keyword        string       `yaml:"keyword" structs:"keyword"`
	RunningSubtext string       `yaml:"running-subtext" structs:"runningsubtext"`
	Title          string       `yaml:"title" structs:"title"`
	WithSpace      bool         `yaml:"with-space" structs:"withspace"`
	Script         ScriptConfig `yaml:"script" structs:"-"`
}

func (s ScriptFilter) ToWorkflowConfig() map[string]interface{} {
	m := structs.Map(s)
	sMap := structs.Map(s.Script)

	for k, v := range sMap {
		m[k] = v
	}

	return m
}

// Clipboard is an object that copies to the clipboard
type Clipboard struct {
	Text string `yaml:"text" structs:"text"`
}

func (c Clipboard) ToWorkflowConfig() map[string]interface{} {
	return structs.Map(c)
}

// OpenURL is an object that opens a URL in a browser
type OpenURL struct{}

func (o OpenURL) ToWorkflowConfig() map[string]interface{} {
	return structs.Map(o)
}

// ScriptConfig is a runnable script in a workflow.
type ScriptConfig struct {
	Content string `yaml:"content" structs:"script"`
	Path    string `yaml:"path" structs:"scriptfile"`
	Type    string `yaml:"type" structs:"-"`
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
	"inline":       8,
}

func (s ScriptConfig) ToWorkflowConfig() map[string]interface{} {
	m := structs.Map(s)
	m["type"] = scriptType[s.Type]
	return m
}
