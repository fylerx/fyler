package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP struct {
		Port int
	}
	CRYPTO struct {
		Passphrase string
	}
	DB struct {
		Host     string
		Username string
		Password string
		Database string
		Port     string
		Sslmode  string
		Timezone string
	}
	S3 struct {
		AccessKeyID     string `mapstructure:"access_key_id"`
		SecretAccessKey string `mapstructure:"secret_access_key"`
		Bucket          string
		Endpoint        string
		Region          string
		DisableSSL      bool
	}
}

var Defaults = map[string]interface{}{
	"http": map[string]string{
		"port": "8080",
	},
	"crypto": map[string]string{
		"passphrase": "passphrasewhichneedstobe32bytes!",
	},
	"db": map[string]string{
		"host":     "postgres",
		"username": "fylerx",
		"password": "fylerx",
		"database": "fylerx",
		"port":     "5432",
		"sslmode":  "disable",
		"timezone": "Europe/Moscow",
	},
	"s3": map[string]interface{}{
		"access_key_id":     "test",
		"secret_access_key": "test",
		"bucket":            "test",
		"endpoint":          "",
		"region":            "",
		"disable_ssl":       true,
	},
}

func Read(appName string, defaults map[string]interface{}, cfg interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetConfigName(appName)
	v.AddConfigPath("./configs/")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	if cfg != nil {
		err := v.Unmarshal(cfg)
		if err != nil {
			return nil, err
		}
	}

	return v, nil
}
