package config

// Script is an Alfred action that runs a script
type Script struct {
	Script ScriptConfig
}

func (s Script) ToWorkflowConfig() map[string]interface{} {
	return s.Script.ToWorkflowConfig()
}
