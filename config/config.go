package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/spf13/viper"
)

var (
	cfg    Config
	doOnce sync.Once
)

type Config struct {
	Application struct {
		Kafka struct {
			Group  string `mapstructure:"GROUP"`
			Broker string `mapstructure:"BROKER"`
			Topic  string `mapstructure:"TOPIC"`
		} `mapstructure:"KAFKA"`
	} `mapstructure:"APPLICATION"`
	Mail struct {
		Host    string `mapstructure:"HOST"`
		Port    int    `mapstructure:"PORT"`
		Address string `mapstructure:"ADDRESS"`
		Sender  string `mapstructure:"SENDER"`
		Pass    string `mapstructure:"PASS"`
	} `mapstructure:"MAIL"`
}

func Get() Config {
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf(fmt.Sprintf("cannot read file env : %v", err))
	}

	doOnce.Do(func() {
		if err := viper.Unmarshal(&cfg); err != nil {
			log.Fatalln("cannot unmarsahl config")
		}
	})

	return cfg
}
