package cloudflare

import (
	"github.com/cloudflare/cloudflare-go"
	"github.com/jbowes/cling"
)

type API struct {
	client *cloudflare.API
}

func NewAPI(email, key string) (*API, error) {
	c, err := cloudflare.New(email, key)
	if err != nil {
		return nil, cling.Wrap(err, "unable to parse CloudFlare credentials")
	}
	return &API{
		client: c,
	}, nil
}
