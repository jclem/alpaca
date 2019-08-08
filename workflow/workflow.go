package workflow

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/jclem/alpaca/config"
)

var argumentType = map[string]int64{
	"required": 0,
	"optional": 1,
	"none":     2,
}

var objectType = map[string]string{
	"keyword":       "alfred.workflow.input.keyword",
	"script":        "alfred.workflow.action.script",
	"script-filter": "alfred.workflow.input.scriptfilter",
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
}

var inlineScript = 8

// Info represents an info.plist in a workflow.
type Info struct {
	BundleID            string                  `plist:"bundleid,omitempty"`
	Connections         map[string][]Connection `plist:"connections,omitempty"`
	CreatedBy           string                  `plist:"createdby,omitempty"`
	Description         string                  `plist:"description,omitempty"`
	Name                string                  `plist:"name,omitempty"`
	Objects             []Object                `plist:"objects,omitempty"`
	Readme              string                  `plist:"readme,omitempty"`
	UIData              map[string]UIData       `plist:"uidata,omitempty"`
	WebAddress          string                  `plist:"webaddress,omitempty"`
	Variables           map[string]string       `plist:"variables,omitempty"`
	VariablesDontExport []string                `plist:"variablesdontexport,omitempty"`
	Version             string                  `plist:"version,omitempty"`
	path                string
}

// NewFromConfig creates a new Info struct from an Alpaca config struct.
func NewFromConfig(path string, c config.Config) (*Info, error) {
	i := Info{
		BundleID:    c.BundleID,
		Connections: map[string][]Connection{},
		CreatedBy:   c.Author,
		Description: c.Description,
		Name:        c.Name,
		Readme:      c.Readme,
		WebAddress:  c.URL,
		Version:     c.Version,
	}

	// Build workflow connections.
	for _, cfgObj := range c.Objects {
		for _, then := range cfgObj.Then {
			conns, ok := i.Connections[cfgObj.UID]
			if !ok {
				i.Connections[cfgObj.UID] = []Connection{}
			}

			// Find the UID for the object we're connecting to.
			var uid string
			for _, cfgObj := range c.Objects {
				if cfgObj.Name == then.Object {
					uid = cfgObj.UID
					break
				}
			}
			if uid == "" {
				return nil, fmt.Errorf("Could not find object %q", then.Object)
			}

			i.Connections[cfgObj.UID] = append(conns, Connection{
				To: uid,
			})
		}
	}

	// Build workflow objects.
	for _, cfgObj := range c.Objects {
		obj := Object{
			Type:    objectType[cfgObj.Type],
			UID:     cfgObj.UID,
			Version: cfgObj.Version,
			Config:  map[string]interface{}{},
		}

		if cfgObj.Script != nil {
			if err := obj.Config.addScriptConfig(path, cfgObj.Script); err != nil {
				return nil, err
			}
		}

		if cfgObj.Type == "keyword" {
			obj.Config["keyword"] = cfgObj.Keyword
			obj.Config["withspace"] = cfgObj.WithSpace
			obj.Config["argumenttype"] = argumentType[cfgObj.Argument]
		}

		i.Objects = append(i.Objects, obj)
	}

	return &i, nil
}

// Connection is a line between two objects.
type Connection struct {
	To              string `plist:"destinationuid,omitempty"`
	Modifiers       int64  `plist:"modifiers,omitempty"`
	ModifierSubtext string `plist:"modifiersubtext,omitempty"`
	VetoClose       bool   `plist:"vitoclose,omitempty"` // NOTE: Yes, "vitoclose"
}

// Object is an object in an Alfred workflow.
type Object struct {
	Config  Config `plist:"config,omitempty"`
	Type    string `plist:"type,omitempty"`
	UID     string `plist:"uid,omitempty"`
	Version int64  `plist:"version,omitempty"`
}

// Config is a generic object configuration object.
type Config map[string]interface{}

func (c *Config) addScriptConfig(path string, script *config.Script) error {
	cfg := *c

	// An inlined path script
	if script.Inline {
		cfg["type"] = scriptType[script.Type]
		scriptPath := filepath.Join(path, script.Path)
		bytes, err := ioutil.ReadFile(scriptPath)
		if err != nil {
			return err
		}
		cfg["script"] = string(bytes)

		return nil
	}

	// An inline script
	if script.Type != "" {
		cfg["type"] = scriptType[script.Type]
		cfg["script"] = script.Content
		return nil
	}

	// A path to a script
	if script.Path != "" {
		cfg["type"] = inlineScript
		cfg["scriptfile"] = script.Path
		return nil
	}

	return nil
}

// UIData represents the position of an object.
type UIData struct {
	XPos int64 `plist:"xpos,omitempty"`
	YPos int64 `plist:"ypos,omitempty"`
}
