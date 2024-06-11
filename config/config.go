package config

import (
	"bufio"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

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
	GO_TESTING_CONTEXT               *testing.T
}

type envAbsPaths struct {
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

func getDefaultConfigValue(fieldName string) string {
	var defaultConfigValue string

	if fieldName == "GO_FIBER_ENVIRONMENT" {
		defaultConfigValue = "development"
	} else if fieldName == "GO_FIBER_SERVER_HOST" {
		defaultConfigValue = "localhost"
	} else if fieldName == "GO_FIBER_SERVER_PORT" {
		defaultConfigValue = "5050"
	} else if fieldName == "GO_FIBER_VAULTS_URL" {
		defaultConfigValue = "http://localhost:8080"
	}

	return defaultConfigValue
}

func scanFileFirstLineToConf(file *os.File, path string, fieldName string, conf *AppConfig,
confElem *reflect.Value) {
	scanner := bufio.NewScanner(file)
	scanner.Scan()

	if contents := scanner.Text(); scanner.Err() != nil {
		log.Fatalf(
			"Error reading contents of '%s' from environment variable %s:\n%s", path, fieldName,
			scanner.Err(),
		)
	} else if fieldName == "GO_FIBER_BEHIND_PROXY" {
		if contents == "true" {
			conf.GO_FIBER_BEHIND_PROXY = true
		} else {
			conf.GO_FIBER_BEHIND_PROXY = false
		}
	} else if fieldName == "GO_FIBER_PROXY_IP_ADDRESSES" {
		conf.GO_FIBER_PROXY_IP_ADDRESSES = strings.Split(contents, ",")
	} else if contents == "" {
		confElem.FieldByName(fieldName).SetString(getDefaultConfigValue(fieldName))
	} else {
		confElem.FieldByName(fieldName).SetString(contents)
	}
}

func loadFileContentsFromPathsToConf(
	conf *AppConfig, pathsType *reflect.Type, pathsValue *reflect.Value, fieldCount int,
) {
	confElem := reflect.ValueOf(conf).Elem()

	for i := 0; i < fieldCount; i++ {
		fieldName := (*pathsType).Field(i).Name
		path := pathsValue.Field(i).Interface().(string)

		if path == "" {
			log.Fatal("Missing or empty environment variable: ", fieldName)
		}

		file, err := os.Open(path)

		if err != nil {
			log.Fatalf(
				"Error opening '%s' from environment variable %s:\n%s", path, fieldName, err,
			)
		}

		defer file.Close()

		scanFileFirstLineToConf(file, path, fieldName, conf, &confElem)
	}
}

func LoadConfigFromEnv(conf *AppConfig) (err error) {
	viper.SetEnvPrefix("GO_FIBER")
	viper.SetConfigFile("./.env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = nil
		} else {
			return
		}
	}

	paths := envAbsPaths{}
	pathsValue := reflect.ValueOf(paths)
	pathsType := pathsValue.Type()
	fieldCount := pathsValue.NumField()

	for i, key, val := 0, "", ""; i < fieldCount; i++ {
		key = pathsType.Field(i).Name

		if val = os.Getenv(key); val != "" {
			viper.Set(key, val)
		}
	}

	if err = viper.Unmarshal(&paths); err != nil {
		return
	} else {
		pathsValue = reflect.ValueOf(paths)
	}

	loadFileContentsFromPathsToConf(conf, &pathsType, &pathsValue, fieldCount)

	return
}
