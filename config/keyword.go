package config

import "github.com/fatih/structs"

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
