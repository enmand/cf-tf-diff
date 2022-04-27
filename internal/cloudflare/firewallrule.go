package cloudflare

import (
	"encoding/json"

	"github.com/cloudflare/cloudflare-go"
	"github.com/jbowes/cling"
)

type FirewallRule struct {
	cloudflare.FirewallRule
}

func (r *FirewallRule) UnmarshalJSON(data []byte) error {
	var f cloudflare.FirewallRule
	if err := json.Unmarshal(data, &f); err != nil {
		return cling.Wrap(err, "failed to unmarshal firewall rule")
	}

	ft := struct {
		FilterID string `json:"filter_id"`
	}{}
	if err := json.Unmarshal(data, &ft); err != nil {
		return cling.Wrap(err, "failed to unmarshal firewall rule contents")
	}

	r.FirewallRule = f
	r.FirewallRule.Filter = cloudflare.Filter{ID: ft.FilterID}

	return nil
}

func (r *FirewallRule) GetFirewallRule(zid string) *Zoned[cloudflare.FirewallRule] {
	return &Zoned[cloudflare.FirewallRule]{
		zoneID: zoneID{zid},
		item:   r.FirewallRule,
	}
}
