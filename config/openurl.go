package config

import "github.com/fatih/structs"

// OpenURL is an object that opens a URL in a browser
type OpenURL struct {
	URL string `yaml:"url" structs:"url"`
}

func (o OpenURL) ToWorkflowConfig() map[string]interface{} {
	return structs.Map(o)
}
