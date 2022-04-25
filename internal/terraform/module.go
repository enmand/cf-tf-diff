package terraform

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/enmand/cf-tf-diff/internal/backends"
	"github.com/enmand/cf-tf-diff/internal/terraform/state"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/jbowes/cling"
)

// ReadBackend reads the .terraform/terraform.tfstate to get a map[string]interface{} of the
// backend configuration from the JSON file.
func GetBackend(p string) (backends.Backend, error) {
	// read the .terraform/terraform.tfstate file
	sf, err := os.Open(path.Join(p, ".terraform/terraform.tfstate"))
	if err != nil {
		return nil, cling.Wrap(err, "unable to read state file")
	}

	// decode the JSON into a map[string]interface{}
	var s state.State
	err = json.NewDecoder(sf).Decode(&s)
	if err != nil {
		return nil, cling.Wrap(err, "unable to decode state file")
	}

	return backends.NewBackend(s.Backend)
}

// GetBackend returns the backend for the given URL.
func GetResources(p string) (map[string]*tfconfig.Resource, error) {
	mod, _ := tfconfig.LoadModule(p)
	if mod.ModuleCalls == nil || len(mod.ModuleCalls) == 0 {
		return nil, cling.Wrap(fmt.Errorf("no backend found in %s", p), "unable to get backend")
	}

	mods, err := getResources(p, mod)
	if err != nil {
		return nil, cling.Wrap(err, "unable to get inner modules")
	}

	_ = mods

	return nil, fmt.Errorf("unable to get backend")
}

// getResources gets all the resources from the given *tf.Module, including any inner modules defined
// by ModuleCalls in the Module recursively.
func getResources(p string, mod *tfconfig.Module) (map[string]*tfconfig.Resource, error) {
	resources := make(map[string]*tfconfig.Resource)

	for mn, m := range mod.ModuleCalls {
		// Get the path to the inner module
		innerPath := path.Join(p, m.Source)

		// Load the inner module
		innerMod, _ := tfconfig.LoadModule(innerPath)

		// Get the resources from the inner module
		innerResources, err := getResources(innerPath, innerMod)
		if err != nil {
			return nil, cling.Wrap(err, "unable to get inner resources")
		}

		// Add the inner resources to the map
		for k, v := range innerResources {
			resources[fmt.Sprintf("module.%s.%s", mn, k)] = v
		}
	}

	// Add the module's resources to the map
	for k, v := range mod.ManagedResources {
		resources[k] = v
	}

	return resources, nil
}
