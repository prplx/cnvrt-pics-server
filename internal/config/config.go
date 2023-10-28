package config

import (
	"flag"
	"os"

	"github.com/prplx/lighter.pics/internal/types"
	"gopkg.in/yaml.v2"
)

func NewConfig(configPath string) (*types.Config, error) {
	config := &types.Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	flag.StringVar(&config.DB.DSN, "db-dsn", "", "Database DSN")
	flag.StringVar(&config.Env, "env", "development", "Environment")
	flag.StringVar(&config.Pusher.AppID, "pusher-app-id", "", "Pusher App ID")
	flag.StringVar(&config.Pusher.Key, "pusher-key", "", "Pusher Key")
	flag.StringVar(&config.Pusher.Secret, "pusher-secret", "", "Pusher Secret")
	flag.StringVar(&config.Pusher.Cluster, "pusher-cluster", "", "Pusher Cluster")
	flag.Parse()

	return config, nil

}
