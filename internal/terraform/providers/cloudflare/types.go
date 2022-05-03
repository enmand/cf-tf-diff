package cloudflare

import (
	"encoding/json"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/jbowes/cling"
)

var (
	// resourceTypes is a map of string names in terraform to struct resource types.
	resourceTypes = map[string]interface{}{
		"certificate_pack":       &CertificatePack{},
		"record":                 &DNSRecord{},
		"filter":                 &Filter{},
		"firewall_rule":          &FirewallRule{},
		"ip_list":                &IPList{},
		"page_rule":              &PageRule{},
		"worker_route":           &WorkerRoute{},
		"worker_script":          &WorkerScript{},
		"zone":                   &Zone{},
		"zone_settings_override": &ZoneSettings{},
	}
)

func FindResourceType(r string) (interface{}, error) {
	if _, ok := resourceTypes[r]; !ok {
		return nil, cling.Errorf("unknown resource type: %s", r)
	}

	return resourceTypes[r], nil
}

// DNSRecord represents a CloudFlare DNS record.
type DNSRecord struct {
	Resource[cloudflare.DNSRecord, struct {
		Hostname string `json:"hostname"`
		Value    string `json:"value"`
	}]
}

func (r *DNSRecord) UnmarshalJSON(data []byte) error {
	err := r.Resource.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	r.Resource.item.ZoneName = r.Resource.parsed.Hostname
	r.Resource.item.Content = r.Resource.parsed.Value

	return nil
}

// CertificatePack represents a CloudFlare certificate pack.
type CertificatePack struct {
	Resource[cloudflare.CertificatePackAdvancedCertificate, interface{}]
}

// Filter represents a CloudFlare filter.
type Filter struct {
	Resource[cloudflare.Filter, interface{}]
}

// FirewallRule represents a CloudFlare firewall rule.
type FirewallRule struct {
	Resource[cloudflare.FirewallRule, struct {
		FilterID string `json:"filter_id"`
	}]
}

func (r *FirewallRule) UnmarshalJSON(data []byte) error {
	err := r.Resource.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	r.Resource.item.Filter = cloudflare.Filter{ID: r.Resource.parsed.FilterID}

	return nil
}

// IPListItem represents an item in an IP list.
type IPListItem struct {
	Resource[cloudflare.IPListItem, struct {
		Value string `json:"value"`
	}]
}

// IPList represents a CloudFlare IP list.
type IPList struct {
	Resource[cloudflare.IPList, interface{}]
	Items []IPListItem `json:"item"`
}

func (r *IPList) UnmarshalJSON(data []byte) error {
	err := r.Resource.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	r.Resource.item.NumItems = len(r.Items)

	return nil
}

// PageRule represents a CloudFlare page rule.
type PageRule struct {
	Resource[cloudflare.PageRule, interface{}]
}

// WorkerRoute represents a CloudFlare worker route.
type WorkerRoute struct {
	Resource[cloudflare.WorkerRoute, struct {
		ScriptName string `json:"script_name"`
	}]
}

func (r *WorkerRoute) UnmarshalJSON(data []byte) error {
	err := r.Resource.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	r.Resource.item.Script = r.Resource.parsed.ScriptName

	return nil
}

// WorkerScript represents a CloudFlare worker script.
type WorkerScript struct {
	Resource[cloudflare.WorkerScript, struct {
		Content string `json:"content"`
		Name    string `json:"name"`
	}]
}

func (r *WorkerScript) UnmarshalJSON(data []byte) error {
	err := r.Resource.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	r.Resource.item.Script = r.Resource.parsed.Content
	r.Resource.item.ID = r.Resource.parsed.Name

	return nil
}

// Zone returns the *cloudflare.Zone associated with this zone, but is able
// to parse strings into booleans, and the plan name into a cloudflare.ZonePlan.
type Zone struct {
	Resource[cloudflare.Zone, interface{}]
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

	z.item = zone
	z.item.Meta = meta
	z.item.Plan = plan

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

// ZoneSettingValues represents the settings values for a ZoneSetting.
type ZoneSettingValues struct {
	InitialSettings       []map[string]interface{} `json:"initial_settings"`
	InitialSettingsReadAt time.Time                `json:"initial_settings_read_at"`
	ReadonlySettings      []string                 `json:"readonly_settings"`
	Settings              []map[string]interface{} `json:"settings"`
}

type ZoneSettings struct {
	Resource[cloudflare.ZoneSetting, ZoneSettingValues] `json:"-"`
	ZoneSettingValues
}

func (r *ZoneSettings) UnmarshalJSON(data []byte) error {
	if err := r.Resource.UnmarshalJSON(data); err != nil {
		return cling.Wrap(err, "failed to unmarshal zone settings")
	}

	r.ZoneSettingValues = r.parsed
	return nil
}
