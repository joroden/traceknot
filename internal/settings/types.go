package settings

type ProviderState struct {
	Binary  string `json:"binary"`
	Enabled bool   `json:"enabled"`
}

type State struct {
	AutostartOn              bool            `json:"autostart_on"`
	RequireWorkItem          bool            `json:"require_work_item"`
	SkipNonInteractivePicker bool            `json:"skip_non_interactive_picker"`
	Hooks                    []ProviderState `json:"hooks"`
	Skills                   []ProviderState `json:"skills"`
}

type Patch struct {
	AutostartOn              *bool           `json:"autostart_on,omitempty"`
	RequireWorkItem          *bool           `json:"require_work_item,omitempty"`
	SkipNonInteractivePicker *bool           `json:"skip_non_interactive_picker,omitempty"`
	Hooks                    map[string]bool `json:"hooks,omitempty"`
	Skills                   map[string]bool `json:"skills,omitempty"`
}
