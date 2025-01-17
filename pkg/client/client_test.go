package client

import (
	"authnet/pkg/config"
	"testing"
)

func Test_ConnectToSandbox(t *testing.T) {
	conf, loadErr := config.LoadConfigFromEnv(false)
	if loadErr != nil {
		t.Fatal(loadErr)
	}
	client := NewAuthNetClient(*conf)
	if _, authErr := client.AuthenticateTest(); authErr != nil {
		t.Fatal(authErr)
	}
}
