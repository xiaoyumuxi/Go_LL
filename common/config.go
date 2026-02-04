package common

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
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

	// 初始解析
	if err := viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Errorf("配置解析失败: %w", err))
	}

	// 开启监听
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("配置文件已修改: %s", e.Name)
		// 重新解析配置
		if err := viper.Unmarshal(&Conf); err != nil {
			log.Printf("配置文件重载失败: %v", err)
		} else {
			log.Printf("配置文件重载成功. 新端口: %d", Conf.Server.Port)
			// 注意：这里虽然更新了 Conf 变量，但已经启动的 Server (gin.Run) 不会自动切换端口。
			// 热更新通常用于动态开关、日志级别、缓存过期时间等运行时可变的参数。
		}
	})
}
