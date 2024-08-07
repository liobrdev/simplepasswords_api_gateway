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
	ENVIRONMENT             string
	BEHIND_PROXY            bool
	PROXY_IP_ADDRESSES      []string
	API_GATEWAY_DB_USER     string
	API_GATEWAY_DB_PASSWORD string
	API_GATEWAY_DB_HOST     string
	API_GATEWAY_DB_PORT     string
	API_GATEWAY_DB_NAME     string
	LOGGER_DB_USER          string
	LOGGER_DB_PASSWORD      string
	LOGGER_DB_HOST          string
	LOGGER_DB_PORT          string
	LOGGER_DB_NAME          string
	SECRET_KEY              string
	API_GATEWAY_HOST				string
	API_GATEWAY_PORT				string
	ADMIN_NAME							string
	ADMIN_EMAIL							string
	ADMIN_PHONE							string
	ADMIN_SALT_1						string
	ADMIN_SALT_2						string
	VAULTS_HOST							string
	VAULTS_PORT							string
	REDIS_PASSWORD          string
	EMAIL_HOST_USER					string
  EMAIL_HOST_PASSWORD			string
	TWILIO_ACCOUNT_SID			string
	TWILIO_AUTH_TOKEN				string
	TWILIO_PHONE_NUMBER			string
	GO_TESTING_CONTEXT			*testing.T
}

type envAbsPaths struct {
	ENVIRONMENT             string
	BEHIND_PROXY            string
	PROXY_IP_ADDRESSES      string
	API_GATEWAY_DB_USER     string
	API_GATEWAY_DB_PASSWORD string
	API_GATEWAY_DB_HOST     string
	API_GATEWAY_DB_PORT     string
	API_GATEWAY_DB_NAME     string
	LOGGER_DB_USER          string
	LOGGER_DB_PASSWORD      string
	LOGGER_DB_HOST          string
	LOGGER_DB_PORT          string
	LOGGER_DB_NAME          string
	REDIS_PASSWORD          string
	SECRET_KEY              string
	API_GATEWAY_HOST				string
	API_GATEWAY_PORT				string
	ADMIN_NAME							string
	ADMIN_EMAIL							string
	ADMIN_PHONE							string
	ADMIN_SALT_1						string
	ADMIN_SALT_2						string
	VAULTS_HOST							string
	VAULTS_PORT							string
	EMAIL_HOST_USER					string
  EMAIL_HOST_PASSWORD			string
	TWILIO_ACCOUNT_SID			string
	TWILIO_AUTH_TOKEN				string
	TWILIO_PHONE_NUMBER			string
}

func getDefaultConfigValue(fieldName string) string {
	var defaultConfigValue string

	if fieldName == "ENVIRONMENT" {
		defaultConfigValue = "development"
	} else if fieldName == "API_GATEWAY_HOST" {
		defaultConfigValue = "localhost"
	} else if fieldName == "API_GATEWAY_PORT" {
		defaultConfigValue = "5050"
	} else if fieldName == "VAULTS_HOST" {
		defaultConfigValue = "localhost"
	} else if fieldName == "VAULTS_PORT" {
		defaultConfigValue = "8080"
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
	} else if fieldName == "BEHIND_PROXY" {
		if contents == "true" {
			conf.BEHIND_PROXY = true
		} else {
			conf.BEHIND_PROXY = false
		}
	} else if fieldName == "PROXY_IP_ADDRESSES" {
		conf.PROXY_IP_ADDRESSES = strings.Split(contents, ",")
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
