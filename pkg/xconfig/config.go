package xconfig

import (
	"fmt"
	"github.com/subosito/gotenv"
	"strings"

	"github.com/spf13/viper"
)

func init() {
	if err := gotenv.Load(); err != nil && !strings.Contains(err.Error(), "no such file or directory") {
		panic(fmt.Errorf("fatal error config .env file: %w", err))
	}
}

func New[T any]() *T {
	// Read config file
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var config T

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("fatal error parsing configuration file: %w", err))
	}

	return &config
}
