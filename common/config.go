package common

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Server     Server     `mapstructure:"server"`
	Datasource Datasource `mapstructure:"datasource"`
	Jwt        Jwt        `mapstructure:"jwt"`
}

type Server struct {
	Port int `mapstructure:"port"`
}

type Datasource struct {
	DriverName string `mapstructure:"driverName"`
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	Database   string `mapstructure:"database"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	Charset    string `mapstructure:"charset"`
}

type Jwt struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"`
}

// 全局配置变量
var Conf *Config

func InitConfig() {
	viper.SetConfigName("config") // 文件名
	viper.SetConfigType("yaml")   // 文件格式
	viper.AddConfigPath(".")      // 在当前目录查找

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("读取配置文件失败: %w", err))
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Errorf("配置解析失败: %w", err))
	}
}
