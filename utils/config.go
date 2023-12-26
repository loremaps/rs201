package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Relays []RelayConfig `mapstructure:"relays"`
}

type RelayConfig struct {
	Group string `mapstructure:"group"`
	Name  string `mapstructure:"name"`
	IP    string `mapstructure:"ip"`
}

func LoadConfig() (config Config, err error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err = viper.Unmarshal(&config)
	return
}

func (c *Config) GetRelayByIP(relayIP string) (relay RelayConfig, err error) {
	for _, r := range c.Relays {
		if r.IP == relayIP {
			relay = r
			return
		}
	}
	err = fmt.Errorf("relay not found")
	return
}
