package conf

import (
	"fileStore/internel/domain"
	"github.com/spf13/viper"
	"sync"
	"time"
)

type Conf struct {
	DbConfig struct {
		Resource string
	}
	RedisConfig struct {
		Addr         string
		Db           int
		Username     string
		Password     string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		DialTimeout  time.Duration
	}
}

var once sync.Once
var config Conf

func init() {
	viper.SetConfigName("Config") // 配置文件名
	viper.SetConfigType("yaml")   // 配置文件类型
	viper.AddConfigPath("./conf") // 配置文件路径

	err := viper.ReadInConfig() // 读取配置文件
	if err != nil {
		panic(err)
	}

}

func GetConfig() Conf {
	dbResourceKey := ""
	//懒汉单例法
	if domain.ServiceName == "user" {
		dbResourceKey = "mysql.user-resource"
	}
	once.Do(func() {
		config = Conf{DbConfig: struct{ Resource string }{Resource: viper.GetString(dbResourceKey)},
			RedisConfig: struct {
				Addr         string
				Db           int
				Username     string
				Password     string
				ReadTimeout  time.Duration
				WriteTimeout time.Duration
				DialTimeout  time.Duration
			}{Addr: viper.GetString("redis.addr"), Db: viper.GetInt("redis.db"), Username: viper.GetString("redis.username"), Password: viper.GetString("redis.password"), ReadTimeout: viper.GetDuration("redis.read_timeout"), WriteTimeout: viper.GetDuration("redis.write_timeout"), DialTimeout: viper.GetDuration("redis.dial_timeout")},
		}
	})
	return config
}
