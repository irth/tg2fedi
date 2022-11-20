package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig(c *Config) error {
	cfgPath, ok := os.LookupEnv("TG2FEDI_CFG")
	if !ok {
		cfgPath = "./tg2fedi.yml"
	}

	cfgF, err := os.Open(cfgPath)
	if err != nil {
		return fmt.Errorf("couldn't open config: %s: %w", cfgPath, err)
	}
	defer cfgF.Close()

	if err = yaml.NewDecoder(cfgF).Decode(&c); err != nil {
		return fmt.Errorf("couldn't decode config: %s: %w", cfgPath, err)
	}

	return nil
}
