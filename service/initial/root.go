package initial

import (
	"fmt"
	"github.com/coinbase/kryptology/service/global"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "gg20_node",
	Short: "GG20_Node is an application node that provides GG20 threshold signature services",
	Run: func(cmd *cobra.Command, args []string) {
		// 读取配置值
		//value := viper.GetInt("stable_window.value")
		//unit := viper.GetString("stable_window.unit")
		//fmt.Printf("Stable Window: %d %s\n", value, unit)
		Run()
	},
}

// Execute 启动命令行应用
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// 定义命令行参数
	rootCmd.Flags().Int("server.port", 8080, "Port to run the server on")
	rootCmd.Flags().String("config", "service/config.yaml", "Port to run the server on")

	if err := viper.BindPFlag("server.port", rootCmd.Flags().Lookup("server.port")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

	// 自动绑定环境变量
	viper.SetEnvPrefix("GG20_NODE")
	viper.AutomaticEnv()

	// 加载配置文件
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// 设置默认配置文件路径
	//viper.AddConfigPath("./service/")
	//viper.SetConfigName("config")
	//viper.SetConfigType("yaml")
	configFile, _ := rootCmd.Flags().GetString("config")
	viper.SetConfigFile(configFile)
	viper.SetDefault("gg20.stableTimeWindow", 5000) // 设置稳定窗口默认值

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Fatal error resources file: %s \n", err.Error())
	}
	if err := viper.Unmarshal(&global.Config); err != nil {
		fmt.Printf("unable to decode into struct %s \n", err.Error())
	}
}
