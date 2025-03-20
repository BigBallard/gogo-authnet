package authnet

var ac *AuthNetClient

func init() {
	conf, loadErr := LoadConfigFromEnv()
	if loadErr != nil {
		panic(loadErr)
	}
	newClient := NewAuthNetClient(*conf)
	ac = &newClient
}
