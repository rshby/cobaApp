package config

import (
	"github.com/spf13/viper"
	"log"
)

type ConfigApp struct {
	App      *App
	Database *Database
	Jaeger   *Jaeger
}

type App struct {
	Name   string `json:"name"`
	Port   int    `json:"port"`
	Author string `json:"author"`
}

type Database struct {
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Name     string `json:"name"`
}

type Jaeger struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}
type Config struct {
	ConfigApp *ConfigApp
}

// method function provider
func NewConfig() IConfig {
	cfg := viper.New()
	cfg.SetConfigFile("config.json")
	cfg.SetConfigType("json")
	cfg.AddConfigPath("./")

	if err := cfg.ReadInConfig(); err != nil {
		log.Fatalf("cant load config : %v", err)
	}

	config := &ConfigApp{
		App: &App{
			Name:   cfg.GetString("app.name"),
			Port:   cfg.GetInt("app.port"),
			Author: cfg.GetString("app.author"),
		},
		Database: &Database{
			Port:     cfg.GetInt("database.port"),
			User:     cfg.GetString("database.user"),
			Password: cfg.GetString("database.password"),
			Host:     cfg.GetString("database.host"),
			Name:     cfg.GetString("database.name"),
		},
		Jaeger: &Jaeger{
			Host: cfg.GetString("jaeger.host"),
			Port: cfg.GetInt("jaeger.port"),
		},
	}
	return &Config{config}
}

func (c *Config) GetConfig() *ConfigApp {
	return c.ConfigApp
}
