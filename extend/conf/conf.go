package conf

import (
	"bytes"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"time"
)

// server conf
type ServerConfig struct {
	RunMode      string        `mapstructure:"runMode"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
	JWTSecret    string        `mapstructure:"jwtSecret"`
	JWTExpire    int           `mapstructure:"jwtExpire"`
	PrefixURL    string        `mapstructure:"PrefixUrl"`
	DBSecret     string        `mapstructure:"dbSecret"`
}

var ServerConf = new(ServerConfig)

// mysql
type DBConfig struct {
	DBType  string `mapstructure:"dbType"`
	ConnStr string `mapstructure:"connStr"`
	Debug   bool   `mapstructure:"debug"`
}

var DBConf = new(DBConfig)

// redis
type RedisConfig struct {
	Host        string        `mapstructure:"host"`
	Port        int           `mapstructure:"port"`
	Password    string        `mapstructure:"password"`
	DBNum       int           `mapstructure:"db"`
	MaxIdle     int           `mapstructure:"maxIdle"`
	MaxActive   int           `mapstructure:"maxActive"`
	IdleTimeout time.Duration `mapstructure:"idleTimeout"`
}

var RedisConf = &RedisConfig{}



func InitConf(runMode string) {
	viper.SetConfigType("YAML")
	if runMode == "debug" {
		data, err := ioutil.ReadFile("config/dev_config.yml")
		if err != nil {
			log.Fatalf("Read 'config/dev_config.yml' fail: %v\n", err)
		}
		_ = viper.ReadConfig(bytes.NewBuffer(data))
		_ = viper.UnmarshalKey("server", ServerConf)
		_ = viper.UnmarshalKey("database", DBConf)
		_ = viper.UnmarshalKey("redis", RedisConf)
	}
	if runMode == "release" {
		data, err := ioutil.ReadFile("config/pro_config.yml")
		if err != nil {
			log.Fatalf("Read 'config/pro_config.yml' fail: %v\n", err)
		}
		_ = viper.ReadConfig(bytes.NewBuffer(data))
		_ = viper.UnmarshalKey("server", ServerConf)
		_ = viper.UnmarshalKey("database", DBConf)
		_ = viper.UnmarshalKey("redis", RedisConf)
	}

}
