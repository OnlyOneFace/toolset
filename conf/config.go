// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/8/13

package conf

import (
	"fmt"
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// @title    LoadConfigFile
// @description   文件配置加载
// @param     configFile        string         "配置文件路径"
func LoadConfigFile(configFile string) {
	configPaths := strings.Split(configFile, ".")
	if len(configPaths) != 2 {
		log.Fatalf("config(%s) error.Only one point is allowed to distinguish extension names,example:xxx.xx\n",
			configFile)
	}
	viper.AddConfigPath("./")
	viper.SetConfigType(configPaths[1])
	viper.SetConfigName(configPaths[0])
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("reading %s config file error.Error:%v\n", configFile, err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config file is changed!")
	})
}
