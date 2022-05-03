package providers

import (
	"strings"

	"github.com/enmand/cf-tf-diff/internal/terraform/providers/cloudflare"
	"github.com/jbowes/cling"
)

type ResourceTypeFunc = func(string) (interface{}, error)

// Providers is a map of provider names to provider functions.
var Providers = map[string]ResourceTypeFunc{
	"cloudflare": cloudflare.FindResourceType,
}

func FindResourceType(full string) (interface{}, error) {
	parts := strings.SplitN(full, "_", 2)
	if len(parts) != 2 {
		return nil, cling.Errorf("invalid resource type: %s", full)
	}

	return FindProviderResourceType(parts[0], parts[1])
}

// FindProvider returns the provider resource for the given provider.
func FindProvider(provider string) (ResourceTypeFunc, error) {
	if p, ok := Providers[provider]; ok {
		return p, nil
	}

	return nil, cling.Errorf("provider %s not found", provider)
}

// FindProviderResource returns a Resource for the given provider and resource.
func FindProviderResourceType(provider, resource string) (interface{}, error) {
	p, err := FindProvider(provider)
	if err != nil {
		return nil, cling.Wrap(err, "failed to find provider")
	}

	r, err := p(resource)
	if err != nil {
		return nil, cling.Wrap(err, "failed to find resource type")
	}

	return r, nil
}
