package oconfig

import (
	"os"
	"strconv"
	"strings"
)

func GetEnvAsString(key string, defaultVal string) string {
	if valStr, exists := os.LookupEnv(key); exists {
		return valStr
	}

	return defaultVal
}

func GetEnvAsInt(name string, defaultVal int) int {
	valStr := GetEnvAsString(name, "")
	if value, err := strconv.Atoi(valStr); err == nil {
		return value
	}

	return defaultVal
}

func GetEnvAsBool(name string, defaultVal bool) bool {
	valStr := GetEnvAsString(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

func GetEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := GetEnvAsString(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}
