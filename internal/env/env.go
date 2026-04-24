package env

import (
	"os"
	"strconv"
	"strings"
)

func GetEnvString(key string, defaultValue string) string {
	value, exist := os.LookupEnv(key)

	if !exist {
		return defaultValue
	}
	return value
}

func GetEnvBool(key string, defaultValue bool) bool {
	value, exist := os.LookupEnv(key)

	if !exist {
		return defaultValue
	}

	castValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return castValue
}

func GetEnvInt(key string, defaultValue int) int {
	value, exist := os.LookupEnv(key)

	if !exist {
		return defaultValue
	}

	castValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return castValue
}

func GetEnvStringSlice(key string, defaultValue []string) []string {
	value, exist := os.LookupEnv(key)

	if !exist {
		return defaultValue
	}
	return strings.Split(value, ",")
}
