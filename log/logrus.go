package log

import (
	"fmt"
	"github.com/sirupsen/logrus"

	"os"
)

var Logger = logrus.New()

func init() {
	// 创建一个 Logrus 实例，将日志输出到文件
	file, err := os.Create("log.txt")
	// 创建一个 Logrus 实例
	Logger.Out = file
	if err != nil {
		fmt.Println("log err :", err)
	}

	// 添加日志记录器
	Logger.Info("This is an info message")

}
