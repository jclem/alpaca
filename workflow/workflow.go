package workflow

import (
	"fmt"

	"github.com/jclem/alpaca/config"
)

// Info represents an info.plist in a workflow.
type Info struct {
	BundleID            string                   `plist:"bundleid,omitempty"`
	Connections         map[string][]Connection  `plist:"connections,omitempty"`
	CreatedBy           string                   `plist:"createdby,omitempty"`
	Description         string                   `plist:"description,omitempty"`
	Name                string                   `plist:"name,omitempty"`
	Objects             []map[string]interface{} `plist:"objects,omitempty"`
	Readme              string                   `plist:"readme,omitempty"`
	UIData              map[string]UIData        `plist:"uidata,omitempty"`
	WebAddress          string                   `plist:"webaddress,omitempty"`
	Variables           map[string]string        `plist:"variables,omitempty"`
	VariablesDontExport []string                 `plist:"variablesdontexport,omitempty"`
	Version             string                   `plist:"version,omitempty"`
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
		Variables:   c.Variables,
	}

	for varName := range i.Variables {
		i.VariablesDontExport = append(i.VariablesDontExport, varName)
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
		obj := cfgObj.ToWorkflowConfig()
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

// UIData represents the position of an object.
type UIData struct {
	XPos int64 `plist:"xpos,omitempty"`
	YPos int64 `plist:"ypos,omitempty"`
}
