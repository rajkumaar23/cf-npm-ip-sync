package internal

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	NPMHost         string        `env:"NPM_HOST"`
	NPMEmail        string        `env:"NPM_EMAIL"`
	NPMPassword     string        `env:"NPM_PASSWORD"`
	NPMAccessListID int64         `env:"NPM_ACCESS_LIST_ID"`
	SyncInterval    time.Duration `env:"SYNC_INTERVAL"`
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
		envValue := os.Getenv(tag)
		if envValue == "" {
			return fmt.Errorf("environment variable %s is required but not set", tag)
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(envValue)
		case reflect.Int64:
			if field.Type() == reflect.TypeOf(time.Duration(0)) {
				durationValue, err := time.ParseDuration(envValue)
				if err != nil {
					return err
				}
				field.SetInt(int64(durationValue))
			} else {
				intValue, err := strconv.Atoi(envValue)
				if err != nil {
					return err
				}
				field.SetInt(int64(intValue))
			}
		}
	}
	return nil
}
