package config

import (
	"github.com/fatih/structs"
)

var alfredMatchMode = map[string]int64{
	"exact-boundary": 0,
	"exact-start":    1,
	"word-match":     2,
}

var argumentTrim = map[string]int64{
	"auto": 0,
	"off":  1,
}

var escaping = map[string]int64{
	"spaces":       1,
	"backquotes":   2,
	"double-quote": 4,
	"brackets":     8,
	"semicolons":   16,
	"dollars":      32,
	"backslashes":  64,
}

var queueMode = map[string]int64{
	"wait":      1,
	"terminate": 2,
}

// ScriptFilter is an Alfred filter that runs a script
type ScriptFilter struct {
	Argument            keywordArgumentType `yaml:"argument" structs:"argumenttype"`
	ArgumentTrim        string              `yaml:"argument-trim" structs:"-"`
	Escaping            []string            `yaml:"escaping" structs:"-"`
	IgnoreEmptyArgument bool                `yaml:"ignore-empty-argument" structs:"-"`
	Keyword             string              `yaml:"keyword" structs:"keyword"`
	RunningSubtitle     string              `yaml:"running-subtitle" structs:"runningsubtext"`
	Subtitle            string              `yaml:"subtitle" structs:"subtext"`
	Title               string              `yaml:"title" structs:"title"`
	WithSpace           bool                `yaml:"with-space" structs:"withspace"`
	Script              ScriptConfig        `yaml:"script" structs:"-"`
	AlfredFilters       *struct {
		Mode string `yaml:"mode"`
	} `yaml:"alfred-filters-results" structs:"-"`
	RunBehavior *struct {
		Immediate  bool   `yaml:"immediate"`
		QueueMode  string `yaml:"queue-mode"`
		QueueDelay string `yaml:"queue-delay"`
	} `yaml:"run-behavior" structs:"-"`
}

var queueDelayCustom = map[string]int64{
	"100ms":  1,
	"200ms":  2,
	"300ms":  3,
	"400ms":  4,
	"500ms":  5,
	"600ms":  6,
	"700ms":  7,
	"800ms":  8,
	"900ms":  9,
	"1000ms": 10,
}

func (s ScriptFilter) ToWorkflowConfig() map[string]interface{} {
	m := structs.Map(s)
	sMap := s.Script.ToWorkflowConfig()
	for k, v := range sMap {
		m[k] = v
	}

	m["argumenttrimmode"] = argumentTrim[s.ArgumentTrim]
	m["argumenttreatemptyqueryasnil"] = s.IgnoreEmptyArgument

	// Filters
	if s.AlfredFilters != nil {
		m["alfredfiltersresults"] = true
		m["alfredfiltersresultsmatchmode"] = alfredMatchMode[s.AlfredFilters.Mode]
	}

	// Escaping
	var escapingSum int64
	for _, esc := range s.Escaping {
		escapingSum = escapingSum + escaping[esc]
	}
	m["escaping"] = escapingSum

	if s.RunBehavior != nil {
		// Run behavior
		m["queuedelayimmediatelyinitially"] = s.RunBehavior.Immediate
		m["queuemode"] = queueMode[s.RunBehavior.QueueMode]

		// Queue Delay
		if s.RunBehavior.QueueDelay == "immediate" {
			m["queuedelaymode"] = 0
		} else if s.RunBehavior.QueueDelay == "automatic" {
			m["queuedelaymode"] = 1
		} else if s.RunBehavior.QueueDelay != "" {
			m["queuedelaymode"] = 2
			m["queuedelaycustom"] = queueDelayCustom[s.RunBehavior.QueueDelay]
		}
	}

	return m
}
