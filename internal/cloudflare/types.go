package cloudflare

import (
	"encoding/json"

	"github.com/cloudflare/cloudflare-go"
	"github.com/jbowes/cling"
)

// State is the State of CloudFlare resources in some JSON state.
type State struct {
	CertificatePacks []*cloudflare.CertificatePack
	Records          []*cloudflare.DNSRecord
	Filters          []*cloudflare.Filter
	FirewallRules    []*cloudflare.FirewallRule
	IPLists          []*cloudflare.IPList
	PageRules        []*cloudflare.PageRule
	WorkerRoutes     []*cloudflare.WorkerRoute
	WorkerScripts    []*cloudflare.WorkerScript
	Zones            []*cloudflare.Zone
	ZoneSettings     []*cloudflare.ZoneSetting
}

// Parse parses a json.RawMessage into a cloudflare API object.
func Parse[T any](data json.RawMessage) (*T, error) {
	var objects T
	if err := json.Unmarshal(data, &objects); err != nil {
		return nil, err
	}
	return &objects, nil
}

type Zone struct {
	cloudflare.Zone
}

func (z *Zone) UnmarshalJSON(data []byte) error {
	var zm map[string]interface{}
	if err := json.Unmarshal(data, &zm); err != nil {
		return cling.Wrap(err, "failed to unmarshal zone meta")
	}

	meta, metaok := zm["meta"].(map[string]interface{})
	delete(zm, "meta")
	delete(zm, "plan")
	data, err := json.Marshal(zm)
	if err != nil {
		return cling.Wrap(err, "failed to marshal zone meta")
	}

	var zone cloudflare.Zone
	if err := json.Unmarshal(data, &zone); err != nil {
		return err
	}
	z.Zone = zone

	plan, ok := zm["plan"].(string)
	if ok {
		z.Plan = cloudflare.ZonePlan{
			ZonePlanCommon: cloudflare.ZonePlanCommon{
				Name: plan,
			},
		}
	}

	if metaok {
		z.Meta = cloudflare.ZoneMeta{}
		prq, ok := meta["page_rule_quota"].(int)
		if ok {
			z.Meta.PageRuleQuota = prq
		}

		wp, ok := meta["wildcard_proxiable"].(string)
		if ok {
			if wp == "true" {
				z.Meta.WildcardProxiable = true
			} else {
				z.Meta.WildcardProxiable = false
			}
		}

		pd, ok := meta["phishing_detected"].(string)
		if ok {
			if pd == "true" {
				z.Meta.PhishingDetected = true
			} else {
				z.Meta.PhishingDetected = false
			}
		}
	}

	return nil
}

func (z *Zone) GetZone() *cloudflare.Zone {
	return &z.Zone
}
