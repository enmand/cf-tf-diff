package backends

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/enmand/cf-tf-diff/internal/terraform/state"
	"github.com/jbowes/cling"
	"github.com/mitchellh/mapstructure"
)

// HTTP represents the "http" backend type for Terraform
type HTTP struct {
	Address       string `mapstructure:"address"`
	LockAddress   string `mapstructure:"lock_address"`
	LockMethod    string `mapstructure:"lock_method"`
	Password      string `mapstructure:"password"`
	UnlockAddress string `mapstructure:"unlock_address"`
	UnlockMethod  string `mapstructure:"unlock_method"`
	Username      string `mapstructure:"username"`
}

// GetStateFile reads the state file from the Address using the given Username and Password.
func (b *HTTP) GetStateFile() (s *state.State, error error) {
	c := &http.Client{
		Timeout: time.Minute,
	}

	_, err := b.sendRequest(c, b.LockMethod, b.LockAddress, nil)
	if err != nil {
		return nil, cling.Wrap(err, "unable to lock state")
	}

	defer func() {
		_, err = b.sendRequest(c, b.UnlockMethod, b.UnlockAddress, nil)
	}()

	resp, err := b.sendRequest(c, "GET", b.Address, nil)
	if err != nil {
		return nil, cling.Wrap(err, "unable to read state file")
	}

	return state.ReadState(resp)
}

// NewHTTPBackend creates a new HTTP backend from the given configuration.
func NewHTTPBackend(b map[string]interface{}) (Backend, error) {
	var h HTTP
	err := mapstructure.Decode(b, &h)
	if err != nil {
		return nil, err
	}

	return &h, nil
}

func (b *HTTP) sendRequest(c *http.Client, method, address string, body []byte) (io.Reader, error) {
	req, err := http.NewRequest(method, address, bytes.NewBuffer(body))
	if err != nil {
		return nil, cling.Wrap(err, "unable to send request")
	}

	if b.Username != "" || b.Password != "" {
		req.SetBasicAuth(b.Username, b.Password)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, cling.Wrap(err, "unable to send request")
	}

	respBody := &bytes.Buffer{}
	_, err = respBody.ReadFrom(resp.Body)
	if err != nil {
		return nil, cling.Wrap(err, "unable to read response body")
	}

	defer resp.Body.Close()

	return respBody, nil
}
