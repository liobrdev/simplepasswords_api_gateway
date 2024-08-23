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
	ADMIN_EMAIL							string
	ADMIN_NAME							string
	ADMIN_PHONE							string
	ADMIN_SALT_1						string
	ADMIN_SALT_2						string
	API_GATEWAY_DB_HOST			string
	API_GATEWAY_DB_NAME			string
	API_GATEWAY_DB_PASSWORD	string
	API_GATEWAY_DB_PORT			string
	API_GATEWAY_DB_USER			string
	API_GATEWAY_HOST				string
	API_GATEWAY_PORT				string
	BEHIND_PROXY						bool
	ENVIRONMENT							string
	LOGGER_DB_HOST					string
	LOGGER_DB_NAME					string
	LOGGER_DB_PASSWORD			string
	LOGGER_DB_PORT					string
	LOGGER_DB_USER					string
	PROXY_IP_ADDRESSES			[]string
	REDIS_PASSWORD					string
	SECRET_KEY							string
	TWILIO_ACCOUNT_SID			string
	TWILIO_AUTH_TOKEN				string
	TWILIO_PHONE_NUMBER			string
	VAULTS_ACCESS_TOKEN			string
	VAULTS_HOST							string
	VAULTS_PORT							string
	GO_TESTING_CONTEXT			*testing.T
}

type envAbsPaths struct {
	ADMIN_EMAIL							string
	ADMIN_NAME							string
	ADMIN_PHONE							string
	ADMIN_SALT_1						string
	ADMIN_SALT_2						string
	API_GATEWAY_DB_HOST			string
	API_GATEWAY_DB_NAME			string
	API_GATEWAY_DB_PASSWORD	string
	API_GATEWAY_DB_PORT			string
	API_GATEWAY_DB_USER			string
	API_GATEWAY_HOST				string
	API_GATEWAY_PORT				string
	BEHIND_PROXY						string
	ENVIRONMENT							string
	LOGGER_DB_HOST					string
	LOGGER_DB_NAME					string
	LOGGER_DB_PASSWORD			string
	LOGGER_DB_PORT					string
	LOGGER_DB_USER					string
	PROXY_IP_ADDRESSES			string
	REDIS_PASSWORD					string
	SECRET_KEY							string
	TWILIO_ACCOUNT_SID			string
	TWILIO_AUTH_TOKEN				string
	TWILIO_PHONE_NUMBER			string
	VAULTS_ACCESS_TOKEN			string
	VAULTS_HOST							string
	VAULTS_PORT							string
}

func scanFileFirstLineToConf(
	file *os.File, conf *AppConfig, confElem *reflect.Value, path, fieldName string, 
) {
	scanner := bufio.NewScanner(file)
	scanner.Scan()

	if contents := scanner.Text(); scanner.Err() != nil {
		log.Fatalf(
			"Error reading contents of '%s' from environment variable %s:\n%s",
			path, fieldName, scanner.Err(),
		)
	} else if contents == "" {
		log.Fatalf("Empty contents of '%s' from environment variable %s", path, fieldName)
	} else if fieldName == "BEHIND_PROXY" {
		if contents == "true" {
			conf.BEHIND_PROXY = true
		} else {
			conf.BEHIND_PROXY = false
		}
	} else if fieldName == "PROXY_IP_ADDRESSES" {
		conf.PROXY_IP_ADDRESSES = strings.Split(contents, ",")
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

		scanFileFirstLineToConf(file, conf, &confElem, path, fieldName)
	}
}

func LoadConfigFromEnv(conf *AppConfig) (err error) {
	viper.SetConfigFile("./.env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		if err.Error() == "open ./.env: no such file or directory" {
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
