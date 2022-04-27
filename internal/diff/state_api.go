package diff

import (
	"fmt"

	cf "github.com/cloudflare/cloudflare-go"
	"github.com/enmand/cf-tf-diff/internal/cloudflare"
	"github.com/enmand/cf-tf-diff/internal/terraform/state"
	"github.com/jbowes/cling"
)

// ParseState parses the *state.State and return a TFCloudflareState using the
// Parse* functions in the cloudflare package.
func ParseState(s *state.State) (*cloudflare.State, error) {
	var tfCloudflareState cloudflare.State

	for _, r := range s.Resources {
		for _, i := range r.Instances {
			switch r.Type {
			case "cloudflare_certificate_pack":
				certificatePack, err := cloudflare.Parse[cf.CertificatePackAdvancedCertificate](i.Attributes)
				if err != nil {
					return nil, cling.Wrap(err, "unable to parse certificate packs")
				}
				tfCloudflareState.CertificatePacks = append(tfCloudflareState.CertificatePacks, certificatePack)
			case "cloudflare_record":
				dnsRecord, err := cloudflare.Parse[cloudflare.DNSRecord](i.Attributes)
				if err != nil {
					return nil, cling.Wrap(err, "unable to parse dns records")
				}
				tfCloudflareState.Records = append(tfCloudflareState.Records, dnsRecord.Get().GetDNSRecord())
			case "cloudflare_filter":
				filter, err := cloudflare.Parse[cf.Filter](i.Attributes)
				if err != nil {
					return nil, cling.Wrap(err, "unable to parse filters")
				}
				tfCloudflareState.Filters = append(tfCloudflareState.Filters, filter)
			case "cloudflare_firewall_rule":
				firewallRule, err := cloudflare.Parse[cloudflare.FirewallRule](i.Attributes)
				if err != nil {
					return nil, cling.Wrap(err, "unable to parse firewall rules")
				}
				tfCloudflareState.FirewallRules = append(
					tfCloudflareState.FirewallRules,
					firewallRule.Get().GetFirewallRule(firewallRule.GetZoneID()),
				)
			case "cloudflare_ip_list":
				ipList, err := cloudflare.Parse[cf.IPList](i.Attributes)
				if err != nil {
					return nil, cling.Wrap(err, "unable to parse ip lists")
				}
				tfCloudflareState.IPLists = append(tfCloudflareState.IPLists, ipList)
			case "cloudflare_page_rule":
				pageRule, err := cloudflare.Parse[cf.PageRule](i.Attributes)
				if err != nil {
					return nil, cling.Wrap(err, "unable to parse page rules")
				}
				tfCloudflareState.PageRules = append(tfCloudflareState.PageRules, pageRule)
			case "cloudflare_worker_route":
				workerRoute, err := cloudflare.Parse[cf.WorkerRoute](i.Attributes)
				if err != nil {
					return nil, cling.Wrap(err, "unable to parse worker routes")
				}
				tfCloudflareState.WorkerRoutes = append(tfCloudflareState.WorkerRoutes, workerRoute)
			case "cloudflare_worker_script":
				workerScript, err := cloudflare.Parse[cf.WorkerScript](i.Attributes)
				if err != nil {
					return nil, cling.Wrap(err, "unable to parse worker scripts")
				}
				tfCloudflareState.WorkerScripts = append(tfCloudflareState.WorkerScripts, workerScript)
			case "cloudflare_zone":
				zone, err := cloudflare.Parse[cloudflare.Zone](i.Attributes)
				if err != nil {
					return nil, cling.Wrap(err, "unable to parse zones")
				}
				tfCloudflareState.Zones = append(tfCloudflareState.Zones, zone.Get().GetZone())
				continue
			case "cloudflare_zone_settings_override":
				zoneSettings, err := cloudflare.Parse[cf.ZoneSetting](i.Attributes)
				if err != nil {
					return nil, cling.Wrap(err, "unable to parse zone settings")
				}
				tfCloudflareState.ZoneSettings = append(tfCloudflareState.ZoneSettings, zoneSettings)
			default:
				return nil, fmt.Errorf("unknown resource type: %s", r.Type)
			}
		}

	}

	return &tfCloudflareState, nil
}
