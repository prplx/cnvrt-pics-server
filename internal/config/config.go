package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/prplx/cnvrt/internal/types"
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

	flag.IntVar(&config.Server.Port, "port", 3001, "Port")
	flag.StringVar(&config.Process.UploadDir, "upload-dir", "./uploads", "Uploads directory")
	flag.StringVar(&config.DB.DSN, "db-dsn", "", "Database DSN")
	flag.StringVar(&config.Env, "env", "development", "Environment")
	flag.StringVar(&config.App.MetricsUser, "metrics-user", "", "Metrics basic auth user")
	flag.StringVar(&config.App.MetricsPassword, "metrics-password", "", "Metrics basic auth password")
	flag.StringVar(&config.Firebase.ProjectID, "firebase-project-id", "", "Firebase project ID")
	flag.StringVar(&config.Server.AllowOrigins, "allow-origins", "", "Allow origins")
	flag.Parse()

	fmt.Printf("%+v\n", config)

	return config, nil
}

func TestConfig() *types.Config {
	return &types.Config{
		Process: struct {
			UploadDir string `yaml:"uploadDir"`
		}{
			UploadDir: "./temp",
		},
		App: struct {
			Name               string `yaml:"name"`
			JobFlushTimeout    int    `yaml:"jobFlushTimeout"`
			MetricsUser        string
			MetricsPassword    string
			SupportedFileTypes string `yaml:"supportedFileTypes"`
			MaxFileCount       int    `yaml:"maxFileCount"`
		}{
			Name:               "cnvrt",
			MaxFileCount:       5,
			SupportedFileTypes: "jpg,jpeg,png",
		}}
}
