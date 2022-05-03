package state

import (
	"encoding/json"
	"io"

	"github.com/enmand/cf-tf-diff/internal/terraform/providers"
	"github.com/jbowes/cling"
)

type InstanceBody interface{}

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

// ResourceHeader represents the header of a resource in the state.
type ResourceHeader struct {
	Mode     string `json:"mode"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
}

type ResourceBody struct {
	Instances []Instance `json:"instances"`
}

// Resource represents a resource in the Terraform state.
type Resource struct {
	ResourceHeader
	ResourceBody
}

func (r *Resource) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &r.ResourceHeader)
	if err != nil {
		return cling.Wrap(err, "failed to unmarshal resource header")
	}

	body := ResourceBody{}
	if err := json.Unmarshal(data, &body); err != nil {
		return cling.Wrap(err, "failed to unmarshal resource instances")
	}

	for _, instance := range body.Instances {
		bt, err := providers.FindResourceType(r.Type)
		if err != nil {
			return cling.Wrap(err, "failed to find resource type")
		}

		if err := json.Unmarshal(instance.Attributes, bt); err != nil {
			return cling.Wrap(err, "failed to unmarshal resource instance")
		}

		instance.Body = bt

		r.ResourceBody.Instances = append(r.ResourceBody.Instances, instance)
	}

	return nil
}

// Instance represents a resource instance in the state.
type Instance struct {
	Attributes   json.RawMessage `json:"attributes"`
	Dependencies []string        `json:"dependencies"`
	Body         InstanceBody    `json:"-"`
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
