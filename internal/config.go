package internal

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	NPMHost         string `env:"NPM_HOST"`
	NPMEmail        string `env:"NPM_EMAIL"`
	NPMPassword     string `env:"NPM_PASSWORD"`
	NPMAccessListID int    `env:"NPM_ACCESS_LIST_ID"`
}

func NewConfig() (*Config, error) {
	c := &Config{}
	if err := c.readFromDotEnv(); err != nil {
		return nil, err
	}
	if err := c.loadVarsIntoConfig(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) readFromDotEnv() error {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error loading .env file: %v", err)
	}
	return nil
}

func (c *Config) loadVarsIntoConfig() error {
	v := reflect.ValueOf(c).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("env")
		if tag == "" {
			continue
		}

		tagParts := strings.Split(tag, ",")
		envKey := tagParts[0]
		isOptional := len(tagParts) > 1 && tagParts[1] == "optional"

		envValue := os.Getenv(envKey)
		if envValue == "" && !isOptional {
			return fmt.Errorf("environment variable %s is required but not set", envKey)
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(envValue)
		case reflect.Int:
			intValue, err := strconv.Atoi(envValue)
			if err != nil {
				return err
			}
			field.SetInt(int64(intValue))
		}
	}
	return nil
}
