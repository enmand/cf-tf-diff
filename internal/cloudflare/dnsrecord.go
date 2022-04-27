package cloudflare

import (
	"encoding/json"

	"github.com/cloudflare/cloudflare-go"
	"github.com/jbowes/cling"
)

// DNSRecord represents a CloudFlare DNS record.
type DNSRecord struct {
	cloudflare.DNSRecord
}

// UnmarshalJSON unmarshals a JSON representation of a DNS record into a DNSRecord.
func (r *DNSRecord) UnmarshalJSON(data []byte) error {
	var dns cloudflare.DNSRecord
	if err := json.Unmarshal(data, &dns); err != nil {
		return cling.Wrap(err, "failed to unmarshal DNS record")
	}
	r.DNSRecord = dns

	rr := struct {
		Hostname string `json:"hostname"`
		Value    string `json:"value"`
	}{}
	if err := json.Unmarshal(data, &rr); err != nil {
		return cling.Wrap(err, "failed to unmarshal DNS record contents")
	}

	r.DNSRecord.ZoneName = rr.Hostname
	r.DNSRecord.Content = rr.Value

	return nil
}

// GetDNSRecords returns the DNS records associated with this DNSRecord
func (r *DNSRecord) GetDNSRecord() *Zoned[cloudflare.DNSRecord] {
	return &Zoned[cloudflare.DNSRecord]{
		zoneID: zoneID{r.ZoneID},
		item:   r.DNSRecord,
	}
}
