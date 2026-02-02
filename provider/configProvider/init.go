package configProvider

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var config *viper.Viper

func Init() {
	fmt.Println("------------------------------------------------------------")
	var err error
	config = viper.New()
	config.SetConfigType("yaml")
	files, err := os.ReadDir(getConfigDir())
	if err != nil {
		log.Fatal(err)
	}
	config.SetConfigName("app")
	config.AddConfigPath("config/")
	_ = config.ReadInConfig()

	for _, f := range files {
		config.SetConfigName(f.Name())
		_ = config.MergeInConfig()
		fmt.Println("Loaded config: " + f.Name())
	}

	if err != nil {
		log.Fatal("error on parsing configuration file", zap.Error(err))
	}
}

func getConfigDir() string {
	configDir := ""
	if dir, err := os.Getwd(); err != nil {
		return ""
	} else {
		configDir = filepath.Join(dir, "config")
	}
	fmt.Println("configDir", configDir)
	return configDir
}

func GetConfig() *viper.Viper {
	return config
}
