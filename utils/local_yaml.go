package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func LoadYaml() *viper.Viper {
	v := viper.New()

	v.SetConfigType("yaml")
	v.SetConfigName("app")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("没有找到app.yaml，请确认在该目录下添加了app.yaml")
			os.Exit(1)
		}
	}
	return v
}