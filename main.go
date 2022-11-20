package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server     string `yaml:"server"`
	ClientName string `yaml:"client-name"`
	Website    string `yaml:"website"`

	path string
}

func LoadConfig(c *Config) error {
	cfgPath, ok := os.LookupEnv("TG2FEDI_CFG")
	if !ok {
		cfgPath = "./tg2fedi.yaml"
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

func main() {
	app := &cli.App{
		Name:  "tg2fedi",
		Usage: "",
		Action: func(*cli.Context) error {
			var config Config
			if err := LoadConfig(&config); err != nil {
				log.Printf("to generate a config, run `%s gen`\n", os.Args[0])
				return err
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
