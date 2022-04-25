package backends

import (
	"os"

	"github.com/enmand/cf-tf-diff/internal/terraform/state"
)

type Local struct {
	file *os.File
}

func (b *Local) GetStateFile() (*state.State, error) {
	return state.ReadState(b.file)
}

func (b *Local) Close() error {
	return b.file.Close()
}

func NewLocalBackend(b map[string]interface{}) (Backend, error) {
	return &Local{}, nil
}
