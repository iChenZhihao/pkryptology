package initial

import (
	"fmt"
	"github.com/coinbase/kryptology/service/global"
)
import "github.com/spf13/viper"

func LoadConfig() {
	viper.AddConfigPath("./service/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Fatal error resources file: %s \n", err.Error())
	}
	if err := viper.Unmarshal(&global.Config); err != nil {
		fmt.Printf("unable to decode into struct %s \n", err.Error())
	}
}
