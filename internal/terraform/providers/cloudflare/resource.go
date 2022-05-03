package cloudflare

import (
	"encoding/json"

	"github.com/jbowes/cling"
)

type zoneID struct {
	ZoneID string `json:"zone_id"`
}

type Resource[T any, P any] struct {
	zoneID

	item   T
	parsed P
}

func (r Resource[T, P]) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &r.item); err != nil {
		return cling.Wrap(err, "failed to unmarshal resource")
	}

	if err := json.Unmarshal(data, &r.parsed); err != nil {
		return cling.Wrap(err, "failed to unmarshal resource contents")
	}

	if err := json.Unmarshal(data, &r.zoneID); err != nil {
		return cling.Wrap(err, "failed to unmarshal resource zone ID")
	}

	return nil
}

func (r Resource[T, P]) Get() T {
	return r.item
}
