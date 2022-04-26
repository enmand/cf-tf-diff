package state

import (
	"encoding/json"
	"io"
)

// State represents the state of a Terraform configuration.
type State struct {
	Version   int        `json:"version"` // TODO: version specific handling
	Backend   Backend    `json:"backend"`
	Modules   []Module   `json:"modules"`
	Resources []Resource `json:"resources"`
}

// Module represents a module in the state.
type Module struct {
	Path      []string            `json:"path"`
	Resources map[string]Resource `json:"resources"`
}

// Resource represents a resource in the Terraform state.
type Resource struct {
	Mode      string     `json:"mode"`
	Type      string     `json:"type"`
	Name      string     `json:"name"`
	Provider  string     `json:"provider"`
	Instances []Instance `json:"instances"`
}

// Instance represents a resource instance in the state.
type Instance struct {
	Attributes   json.RawMessage `json:"attributes"`
	Dependencies []string        `json:"dependencies"`
}

// Primary represents the primary instance of a resource in the state.
type Primary struct {
	ID         string          `json:"id"`
	Attributes json.RawMessage `json:"attributes"`
}

// Backend represents the backend configuration for the state.
type Backend struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

// ReadState reads a State from the given reader.
func ReadState(br io.Reader) (*State, error) {
	var s State
	if err := json.NewDecoder(br).Decode(&s); err != nil {
		return nil, err
	}
	return &s, nil
}
