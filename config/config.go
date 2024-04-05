package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/chirag1807/logease"
	"github.com/chirag1807/task-management-system/api/model/dto"
	"github.com/spf13/viper"
)

var Config dto.Config
var JWtSecretKey dto.JWTSecret
var LoggerInstance logease.SlogLoggerInstance

// LoadConfig uses viper package to load all env variables into config(package:dto) struct via above declared global Config variable.
// Moreover it also read secret.json file of .config directory and load content into above declared global JWtSecretKey variable.
func LoadConfig(envFilePath string, secretJsonFilePath string) {
	viper.AutomaticEnv()
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath(envFilePath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatal(err)
	}

	jwtSecretKeyFileContent, err := os.ReadFile(secretJsonFilePath)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(jwtSecretKeyFileContent, &JWtSecretKey); err != nil {
		log.Fatal(err)
	}

	loggerInstance, err := logease.InitLogease(true, Config.TeamsWebHookURL, logease.Slog)
	if err != nil {
		log.Fatal(err)
	}
	LoggerInstance = loggerInstance.(logease.SlogLoggerInstance)
}
