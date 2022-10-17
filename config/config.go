package config

import (
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

type AppConfig struct {
	GO_FIBER_ENVIRONMENT             string
	GO_FIBER_BEHIND_PROXY            bool
	GO_FIBER_PROXY_IP_ADDRESSES      []string
	GO_FIBER_API_GATEWAY_DB_USER     string
	GO_FIBER_API_GATEWAY_DB_PASSWORD string
	GO_FIBER_API_GATEWAY_DB_HOST     string
	GO_FIBER_API_GATEWAY_DB_PORT     string
	GO_FIBER_API_GATEWAY_DB_NAME     string
	GO_FIBER_LOGGER_DB_USER          string
	GO_FIBER_LOGGER_DB_PASSWORD      string
	GO_FIBER_LOGGER_DB_HOST          string
	GO_FIBER_LOGGER_DB_PORT          string
	GO_FIBER_LOGGER_DB_NAME          string
	GO_FIBER_REDIS_PASSWORD          string
	GO_FIBER_SECRET_KEY              string
	GO_FIBER_SERVER_HOST             string
	GO_FIBER_SERVER_PORT             string
	GO_FIBER_VAULTS_URL              string
}

type envFileAbsPaths struct {
	GO_FIBER_ENVIRONMENT             string
	GO_FIBER_BEHIND_PROXY            string
	GO_FIBER_PROXY_IP_ADDRESSES      string
	GO_FIBER_API_GATEWAY_DB_USER     string
	GO_FIBER_API_GATEWAY_DB_PASSWORD string
	GO_FIBER_API_GATEWAY_DB_HOST     string
	GO_FIBER_API_GATEWAY_DB_PORT     string
	GO_FIBER_API_GATEWAY_DB_NAME     string
	GO_FIBER_LOGGER_DB_USER          string
	GO_FIBER_LOGGER_DB_PASSWORD      string
	GO_FIBER_LOGGER_DB_HOST          string
	GO_FIBER_LOGGER_DB_PORT          string
	GO_FIBER_LOGGER_DB_NAME          string
	GO_FIBER_REDIS_PASSWORD          string
	GO_FIBER_SECRET_KEY              string
	GO_FIBER_SERVER_HOST             string
	GO_FIBER_SERVER_PORT             string
	GO_FIBER_VAULTS_URL              string
}

func (paths envFileAbsPaths) loadFileContentsToConf(conf *AppConfig) {
	pathsValue := reflect.ValueOf(paths)
	pathsType := pathsValue.Type()
	confElem := reflect.ValueOf(conf).Elem()

	for i, fieldCount := 0, pathsValue.NumField(); i < fieldCount; i++ {
		fieldName := pathsType.Field(i).Name
		path := pathsValue.Field(i).Interface().(string)

		if path == "" {
			log.Fatalln("Missing or empty environment variable:", fieldName)
		}

		if contents, err := os.ReadFile(path); err != nil {
			log.Fatalf(
				"Error reading contents of '%s' from variable %s in '.env' file:\n%s",
				path, fieldName, err,
			)
		} else if fieldName == "GO_FIBER_BEHIND_PROXY" {
			if string(contents) == "true" {
				conf.GO_FIBER_BEHIND_PROXY = true
			} else {
				conf.GO_FIBER_BEHIND_PROXY = false
			}
		} else if fieldName == "GO_FIBER_PROXY_IP_ADDRESSES" {
			conf.GO_FIBER_PROXY_IP_ADDRESSES = strings.Split(string(contents), ",")
		} else {
			confElem.FieldByName(fieldName).SetString(string(contents))
		}
	}
}

func LoadConfigFromEnvFile(conf *AppConfig) (err error) {
	viper.SetConfigFile("./.env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	var paths envFileAbsPaths

	if err = viper.Unmarshal(&paths); err != nil {
		return
	}

	paths.loadFileContentsToConf(conf)

	return
}
