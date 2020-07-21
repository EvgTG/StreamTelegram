package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

//Config Struct
type Config struct {
	viper *viper.Viper
}

//New init Config struct
func New() *Config {
	viper := viper.New()
	viper.AutomaticEnv()

	if configName := viper.GetString("CONFIG_NAME"); configName != "" {
		godotenv.Load(configName)
	}

	return &Config{viper: viper}
}

//GetString gets dconfig value
func (c *Config) GetString(key string) string {
	return c.viper.GetString(key)
}

//IsSet check if key is exists in config
func (c *Config) IsSet(key string) bool {
	return c.viper.IsSet(key)
}

/*
Get(key string) : interface{}
GetBool(key string) : bool
GetFloat64(key string) : float64
GetInt(key string) : int
GetIntSlice(key string) : []int
GetString(key string) : string
GetStringMap(key string) : map[string]interface{}
GetStringMapString(key string) : map[string]string
GetStringSlice(key string) : []string
GetTime(key string) : time.Time
GetDuration(key string) : time.Duration
IsSet(key string) : bool

*/

// author github.com/rmukhamet/
