package api

import (
	"github.com/valkey-io/valkey-go"
)

var (
	hasInit      bool
	valkeyClient valkey.Client
)

const (
	PostTtl         = 5 * 60       // seconds
	TagTtl          = 24 * 60 * 60 // seconds
	TagSearchTtl    = 60 * 60      // seconds
	SettingShareTtl = 24 * 60 * 60 // seconds
)

func InitValkey() error {
	if hasInit {
		return nil
	}

	vc, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{ValkeyAddr}})
	if err != nil {
		return err
	}

	valkeyClient = vc
	hasInit = true

	return nil
}

func Valkey() valkey.Client {
	return valkeyClient
}
