package backends

import (
	"fmt"

	"github.com/enmand/cf-tf-diff/internal/terraform/state"
)

// BackendType is the Type of the backend
type BackendType string

var (
	// HTTPBackend is the "http" backend type for Terraform
	HTTPBackend BackendType = "http"
	// S3Backend is the "s3" backend type for Terraform
	S3Backend BackendType = "s3"

	// LocalBackend is the "local" backend type for Terraform
	LocalBackend BackendType = "local"

	// BackendMap is a map of backend types to their respective backend
	// implementations.
	BackendMap = map[BackendType]func(map[string]interface{}) (Backend, error){
		HTTPBackend:  NewHTTPBackend,
		S3Backend:    NewS3Backend,
		LocalBackend: NewLocalBackend,
	}
)

// NewBackend creates a new backend from the given configuration.
func NewBackend(b state.Backend) (Backend, error) {
	if b.Type == "" {
		return nil, fmt.Errorf("backend type is required")
	}

	f, ok := BackendMap[BackendType(b.Type)]
	if !ok {
		return nil, fmt.Errorf("unknown backend type: %s", b.Type)
	}

	return f(b.Config)
}

// Backend represents a Terraform backend.
type Backend interface {
	GetStateFile() (*state.State, error)
}
