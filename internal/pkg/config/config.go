package config

import (
	"sync"

	"github.com/spf13/viper"
)

var (
	cfg  *viper.Viper
	once sync.Once
)

func Load() *viper.Viper {
	once.Do(func() {
		cfg = viper.New()
		cfg.SetConfigName("config")
		cfg.AddConfigPath("./configs")
		cfg.SetConfigType("yaml")
		cfg.AutomaticEnv()
		_ = cfg.ReadInConfig() // 忽略找不到文件错误
	})
	return cfg
}
