package gogo_authnet

import (
	"testing"
)

func Test_ConnectToSandbox(t *testing.T) {
	conf, loadErr := LoadConfigFromEnv(false)
	if loadErr != nil {
		t.Fatal(loadErr)
	}
	client := NewAuthNetClient(*conf)
	if _, authErr := client.AuthenticateTest(); authErr != nil {
		t.Fatal(authErr)
	}
}
