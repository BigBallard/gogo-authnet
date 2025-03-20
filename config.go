package authnet

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	EnvGoGoAuthnetConfig = "GOGO_AUTHNET_CONFIG"
)

// Config is the configuration for the API. Refer to the [config documentation] for key references.
//
// [config documentation]: https://github.com/BigBallard/gogo-authnet/blob/master/docs/CONFIG.md
type Config struct {
	AuthnetHost string `json:"authnet-host,omitempty"`
	Auth        *Auth  `json:"auth,omitempty"` // API authorization credentials
}

// Auth provides API authorization credentials
type Auth struct {
	ApiLoginId     string `json:"api-login-id,omitempty"`
	TransactionKey string `json:"transaction-key,omitempty"`
}

// aggregate applies the configuration values set through the CLI and the environment variables.
func aggregate(config *Config) {
	arguments := os.Args[1:]
	for i := 0; i < len(arguments); i++ {
		arg := arguments[i]
		// Not an argument
		if !strings.HasPrefix(arg, "-") {
			continue
		}

		value := arguments[i+1]
		i++
		switch arg {
		case "-AUTHNET_HOST":
			config.AuthnetHost = value
		case "-AUTH_API_LOGIN_ID":
			if config.Auth == nil {
				config.Auth = new(Auth)
			}
			config.Auth.ApiLoginId = value
		case "-AUTH_TRANSACTION_KEY":
			if config.Auth == nil {
				config.Auth = new(Auth)
			}
			config.Auth.TransactionKey = value
		}
	}

	if value, ok := os.LookupEnv("AUTHNET_HOST"); ok {
		config.AuthnetHost = value
	}
	if value, ok := os.LookupEnv("AUTH_API_LOGIN_ID"); ok {
		if config.Auth == nil {
			config.Auth = new(Auth)
		}
		config.Auth.ApiLoginId = value
	}
	if value, ok := os.LookupEnv("AUTH_TRANSACTION_KEY"); ok {
		if config.Auth == nil {
			config.Auth = new(Auth)
		}
		config.Auth.TransactionKey = value
	}
}

// LoadConfigFromFile attempts to load the config JSON file from the path provided. Aggregate determines if the config will
// be aggregated with environment variables. Any entry in the config file can be provided from the environment as needed
// by following the naming convention:
func LoadConfigFromFile(path string) (*Config, error) {
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

	aggregate(&config)

	if config.Auth == nil {
		return nil, errors.New("no auth config found")
	}
	return &config, nil
}

// LoadConfigFromEnv checks for the environment variable GOGO_AUTHNET_CONFIG which should be a full system path to a
// JSON file that conforms the [config.Config] schema.
func LoadConfigFromEnv() (*Config, error) {
	configPath, found := os.LookupEnv(EnvGoGoAuthnetConfig)
	if !found {
		return nil, fmt.Errorf("environment variable %s not set", EnvGoGoAuthnetConfig)
	}
	return LoadConfigFromFile(configPath)
}

// LoadConfig creates the configuration from values set through CLI arguments and the environment variables. The same
// validation is conducted as if loaded from a file.
func LoadConfig() (*Config, error) {
	var config Config
	aggregate(&config)
	if config.Auth == nil {
		return nil, errors.New("no auth config found")
	}
	return &config, nil
}
