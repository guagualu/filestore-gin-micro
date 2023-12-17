package conf

import "github.com/spf13/viper"

func init() {
	viper.SetConfigName("config") // 配置文件名
	viper.SetConfigType("yaml")   // 配置文件类型
	viper.AddConfigPath(".")      // 配置文件路径

	err := viper.ReadInConfig() // 读取配置文件
	if err != nil {
		panic(err)
	}

}
