package client

import (
	config2 "authnet/pkg/config"
	"testing"
)

func Test_ConnectToSandbox(t *testing.T) {
	config, loadErr := config2.LoadConfigFromEnv(false)
	if loadErr != nil {
		t.Fatal(loadErr)
	}
	client := NewAuthNetClient(*config)
	if _, authErr := client.AuthenticateTest(); authErr != nil {
		t.Fatal(authErr)
	}
}
