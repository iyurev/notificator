package sender

import (
	"testing"
)

func TestNewGlobalConfig(t *testing.T) {
	config, err := NewGlobalConfig()
	if err != nil {
		t.Fatal(err)
	}
	_ = config
	t.Logf("Global config: %v\n", globalConfig)
	r, err := globalConfig.GetProjectRecipient("myapp2")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("myapp1 recipient data: %d\n", r.ChatID)
}
