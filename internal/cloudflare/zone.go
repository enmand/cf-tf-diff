package cloudflare

import (
	"encoding/json"

	"github.com/cloudflare/cloudflare-go"
	"github.com/jbowes/cling"
)

// Zone returns the *cloudflare.Zone associated with this zone, but is able
// to parse strings into booleans, and the plan name into a cloudflare.ZonePlan.
type Zone struct {
	cloudflare.Zone
}

// UnmarshalJSON unmarshals a JSON representation of a zone into a Zone.
func (z *Zone) UnmarshalJSON(data []byte) error {
	var zm map[string]interface{}
	if err := json.Unmarshal(data, &zm); err != nil {
		return cling.Wrap(err, "failed to unmarshal zone meta")
	}

	meta := cleanZoneMeta(zm)
	plan := stringToPlan(zm)

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
	z.Zone.Meta = meta
	z.Zone.Plan = plan

	return nil
}

func stringToPlan(zm map[string]interface{}) cloudflare.ZonePlan {
	var plan cloudflare.ZonePlan
	p, ok := zm["plan"].(string)
	if ok {
		plan.ZonePlanCommon = cloudflare.ZonePlanCommon{
			Name: p,
		}
	}

	return plan
}

func cleanZoneMeta(zm map[string]interface{}) cloudflare.ZoneMeta {
	var meta cloudflare.ZoneMeta
	mm, ok := zm["meta"].(map[string]interface{})
	if ok {
		prq, ok := mm["page_rule_quota"].(int)
		if ok {
			meta.PageRuleQuota = prq
		}

		wp, ok := mm["wildcard_proxiable"].(string)
		if ok {
			if wp == "true" {
				meta.WildcardProxiable = true
			} else {
				meta.WildcardProxiable = false
			}
		}

		pd, ok := mm["phishing_detected"].(string)
		if ok {
			if pd == "true" {
				meta.PhishingDetected = true
			} else {
				meta.PhishingDetected = false
			}
		}
	}

	return meta
}

// GetZone returns the *cloudflare.Zone associated with this zone.
func (z *Zone) GetZone() *Zoned[cloudflare.Zone] {
	return &Zoned[cloudflare.Zone]{
		zoneID: zoneID{z.ID},
		item:   z.Zone,
	}
}
