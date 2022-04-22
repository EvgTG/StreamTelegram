package config

import (
	"StreamTelegram/go-log"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	viper *viper.Viper
}

func New() *Config {
	viper := viper.New()
	viper.SetConfigName("cfg")
	viper.AddConfigPath("files/")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(errors.Wrap(err, "cfg viper.ReadInConfig()"))
	}

	return &Config{viper: viper}
}

func (c *Config) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

func (c *Config) GetString(key string) string {
	return c.viper.GetString(key)
}

func (c *Config) GetInt64(key string) int64 {
	return c.viper.GetInt64(key)
}

func (c *Config) GetInt(key string) int {
	return c.viper.GetInt(key)
}

func (c *Config) GetIntSlice64(key string) []int64 {
	// viper.GetIntSlice почему-то не работает, хочу разделять запятыми
	sliceStr := c.viper.GetString(key)
	slice64 := make([]int64, 0)

	for _, valueStr := range strings.Split(sliceStr, ",") {
		value, err := strconv.ParseInt(valueStr, 10, 64)
		if err == nil {
			slice64 = append(slice64, value)
		}
	}

	return slice64
}

func (c *Config) GetTimeLocation(key string) (loc *time.Location, err error) {
	locStr := c.viper.GetString(key)

	if locStr == "" {
		loc, _ = time.LoadLocation("")
	} else {
		loc, err = time.LoadLocation(locStr)
		if err != nil {
			return nil, errors.Wrap(err, "time.LoadLocation")
		}
	}

	return
}

func (c *Config) IsSet(key string) bool {
	return c.viper.IsSet(key)
}

/*
func (c *Config) (key string)  {
	return
}
*/
