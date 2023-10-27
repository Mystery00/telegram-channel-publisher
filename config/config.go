package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func init() {
	configHome, exist := os.LookupEnv(EnvConfigHome)
	if !exist {
		//不存在环境变量，尝试取可执行文件目录下的config目录
		configHome = "etc"
	}
	//判断配置文件目录是否存在
	if !exists(configHome) {
		//配置文件目录不存在，拒绝启动
		logrus.Fatal("config file not exist")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configHome)
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		logrus.Info("config file updated")
		err := viper.ReadInConfig()
		if err != nil {
			logrus.Fatal("read config failed", err)
		}
	})
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatal("read config failed", err)
	}
}
