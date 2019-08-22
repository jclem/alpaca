package config

import (
	"github.com/fatih/structs"
	yaml "gopkg.in/yaml.v3"
)

type keywordArgumentType string

var (
	keywordArgumentRequired keywordArgumentType = "required"
	keywordArgumentOptional keywordArgumentType = "optional"
	keywordArgumentNone     keywordArgumentType = "none"
)

var argumentType = map[keywordArgumentType]int64{
	"required": 0,
	"optional": 1,
	"none":     2,
}

// Keyword is an object triggered by a keyword
type Keyword struct {
	Keyword   string              `yaml:"keyword" structs:"keyword"`
	WithSpace bool                `yaml:"with-space" structs:"withspace"`
	Title     string              `yaml:"title" structs:"text"`
	Subtitle  string              `yaml:"subtitle" structs:"subtext"`
	Argument  keywordArgumentType `yaml:"argument" structs:"argumenttype"`
}

func (k *Keyword) UnmarshalYAML(node *yaml.Node) error {
	type alias Keyword
	as := alias{WithSpace: true}
	if err := node.Decode(&as); err != nil {
		return err
	}

	*k = Keyword(as)

	return nil
}

func (k Keyword) ToWorkflowConfig() map[string]interface{} {
	m := structs.Map(k)
	m["argumenttype"] = argumentType[k.Argument]
	return m
}
