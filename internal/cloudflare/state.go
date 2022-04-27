package cloudflare

import (
	"encoding/json"

	"github.com/cloudflare/cloudflare-go"
	"github.com/jbowes/cling"
)

type zoneID struct {
	ZoneID string `json:"zone_id"`
}

// Zoned represents a CloudFlare API object with a ZoneID.
type Zoned[T any] struct {
	zoneID
	item T
}

// GetZoneID returns the ZoneID of a CloudFlare API object.
func (z Zoned[T]) GetZoneID() string {
	return z.ZoneID
}

// Get returns the CloudFlare API object.
func (z Zoned[T]) Get() *T {
	return &z.item
}

// State is the State of CloudFlare resources in some JSON state.
type State struct {
	CertificatePacks []*Zoned[cloudflare.CertificatePackAdvancedCertificate]
	Records          []*Zoned[cloudflare.DNSRecord]
	Filters          []*Zoned[cloudflare.Filter]
	FirewallRules    []*Zoned[cloudflare.FirewallRule]
	IPLists          []*Zoned[cloudflare.IPList]
	PageRules        []*Zoned[cloudflare.PageRule]
	WorkerRoutes     []*Zoned[cloudflare.WorkerRoute]
	WorkerScripts    []*Zoned[cloudflare.WorkerScript]
	Zones            []*Zoned[cloudflare.Zone]
	ZoneSettings     []*Zoned[cloudflare.ZoneSetting]
}

// Parse parses a json.RawMessage into a cloudflare API object.
func Parse[T any](data json.RawMessage) (*Zoned[T], error) {
	var objects Zoned[T]
	if err := json.Unmarshal(data, &objects.item); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &objects.zoneID); err != nil {
		return nil, cling.Wrap(err, "failed to unmarshal zone ID")
	}

	return &objects, nil
}
