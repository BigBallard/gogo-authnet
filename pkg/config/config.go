package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const (
	EnvGoGoAuthnetConfig = "GOGO_AUTHNET_CONFIG"
)

// Config is the configuration for the API. Refer to the [scheme documentation].
//
// [scheme documentation]: https://github.com/BigBallard/gogo-authnet.git
type Config struct {
	AuthnetHost string `json:"authnet-host,omitempty"`
	Auth        *Auth  `json:"auth,omitempty"` // API authorization credentials
}

// Auth provides API authorization credentials
type Auth struct {
	ApiLoginId    string `json:"api-login-id"`
	TransactionId string `json:"transaction-id"`
}

// LoadConfigFile attempts to load the config JSON file from the path provided. Aggregate determines if the config will
// be aggregated with environment variables. Any entry in the config file can be provided from the environment as needed
// by following the naming convention:
func LoadConfigFile(path string, aggregate bool) (*Config, error) {
	bytes, readErr := os.ReadFile(path)
	if readErr != nil {
		return nil, readErr
	}

	var config Config

	if unmarshallErr := json.Unmarshal(bytes, &config); unmarshallErr != nil {
		return nil, errors.Join(errors.New("invalid config file format"), unmarshallErr)
	}

	if len(config.AuthnetHost) == 0 {
		config.AuthnetHost = "https://api.authorize.net"
	}
	if config.Auth == nil {
		return nil, errors.New("no auth config found")
	}
	return &config, nil
}

// LoadConfigFromEnv checks for the environment variable GOGO_AUTHNET_CONFIG which should be a full system path to a
// JSON file that conforms the [config.Config] schema.
//
// Aggregate determines if the config will be aggregated with
// other environment variables.
func LoadConfigFromEnv(aggregate bool) (*Config, error) {
	configPath, found := os.LookupEnv(EnvGoGoAuthnetConfig)
	if !found {
		return nil, fmt.Errorf("environment variable %s not set", EnvGoGoAuthnetConfig)
	}
	return LoadConfigFile(configPath, aggregate)
}
