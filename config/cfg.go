package config

import (
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type (
	Config struct {
		API        API        `mapstructure:"api"`
		Logger     Logger     `mapstructure:"logger"`
		DB         DB         `mapstructure:"db"`
		FileServer FileServer `mapstructure:"file_server"`
		Invoice    Invoice    `mapstructure:"invoice"`
		Services   []string   `mapstructure:"services"`
		JWT        JWT        `mapstructure:"jwt"`
	}

	JWT struct {
		KeysPath string        `mapstructure:"keys_path"`
		TTL      time.Duration `mapstructure:"ttl"`
	}

	API struct {
		Port     int `mapstructure:"port"`
		WTimeout int `mapstructure:"write_timeout"`
		RTimeout int `mapstructure:"read_timeout"`
		GTimeout int `mapstructure:"graceful_timeout"`
	}

	Invoice struct {
		Frequency    time.Duration `mapstructure:"freq"`
		GenerateTime string        `mapstructure:"generate_time"`
	}

	FileServer struct {
		Endpoint  string `mapstructure:"endpoint"`
		Port      string `mapstructure:"port"`
		AccessKey string `mapstructure:"access_key"`
		SecretKey string `mapstructure:"secret_key"`
		SSL       bool   `mapstructure:"ssl"`
	}

	Logger struct {
		Encoding    string   `mapstructure:"encoding"`
		OutputPaths []string `mapstructure:"output_paths"`
	}

	DB struct {
		Provider        string `mapstructure:"provider"`
		Host            string `mapstructure:"host"`
		Port            int    `mapstructure:"port"`
		DB              string `mapstructure:"db_name"`
		User            string `mapstructure:"user"`
		Password        string `mapstructure:"password"`
		SSL             string `mapstructure:"ssl"`
		Timeout         string `mapstructure:"timeout"`
		Retry           int    `mapstructure:"retry"`
		RandomDataCount int    `mapstructure:"random_data_count"`
	}
)

func LoadConfig() (*Config, error) {
	var config *Config
	viper.AutomaticEnv()
	viper.SetDefault("PETSTORE_CONFIG", "config.yaml")
	viper.SetConfigFile(viper.GetString("PETSTORE_CONFIG"))
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "unable to read config with filepath")
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal config to struct")
	}

	catchParamsFromEnv(config)

	return config, nil
}

func catchParamsFromEnv(config *Config) *Config {
	if viper.IsSet("POSTGRES_HOST") {
		config.DB.Host = viper.GetString("POSTGRES_HOST")
	}

	if viper.IsSet("MINIO_HOST") {
		config.FileServer.Endpoint = viper.GetString("MINIO_HOST")
	}

	if viper.IsSet("MINIO_ACCESS_KEY") {
		config.FileServer.AccessKey = viper.GetString("MINIO_ACCESS_KEY")
	}

	if viper.IsSet("MINIO_SECRET_KEY") {
		config.FileServer.SecretKey = viper.GetString("MINIO_SECRET_KEY")
	}

	if viper.IsSet("MINIO_PORT") {
		config.FileServer.Port = viper.GetString("MINIO_PORT")
	}

	return config
}
