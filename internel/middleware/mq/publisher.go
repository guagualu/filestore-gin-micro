package mq

import (
	"fileStore/conf"
	"fmt"
	"sync"

	"github.com/streadway/amqp"
)

var (
	rabconn    *amqp.Connection
	rabchannel *amqp.Channel
	once       sync.Once
)

type MqFileInfo struct {
	FileHash    string
	FileName    string
	CurLocateAt string
}

//
//func initMq() {
//	var err error
//	config := conf.GetConfig()
//	//1、获取mq连接
//	Rabconn, err = amqp.Dial(config.MqAddr)
//	if err != nil {
//		fmt.Println("rabconn dail err:", err)
//		return
//	}
//	//2、获取mq的channel
//	rabchannel, err = Rabconn.Channel()
//	if err != nil {
//		fmt.Println("rabconn channel err:", err)
//		return
//	}
//}

//func GetRabconn() *amqp.Connection {
//	once.Do(func() {
//		var err error
//		config := conf.GetConfig()
//		//1、获取mq连接
//		rabconn, err = amqp.Dial(config.MqAddr)
//		if err != nil {
//			fmt.Println("rabconn dail err:", err)
//			return
//		}
//	})
//	return rabconn
//}

func GetRabchannel() *amqp.Channel {
	once.Do(func() {
		var err error
		config := conf.GetConfig()
		//1、获取mq连接
		rabconn, err = amqp.Dial(config.MqAddr)
		if err != nil {
			fmt.Println("rabconn dail err:", err)
			return
		}
		//2、获取mq的channel
		rabchannel, _ = rabconn.Channel()
	})
	return rabchannel
}

func Rabpublish(routekey, msg string) {
	//如果在web控制台创建了exchange 1可以不用  发布端只负责发布到交换机
	config := conf.GetConfig()
	//1、先定义出exchange
	rabchannel := GetRabchannel()
	err := rabchannel.ExchangeDeclare( //??? 属性细看
		config.MqUploadExchangeName,
		"direct", //路由类型
		true,     //是否持久化
		false,    //是否丢失自动删除
		false,    //是否具有排他性
		false,    //是否堵塞
		nil,      //额外属性
	)
	if err != nil {
		fmt.Println("rabconn channel err:", err)
		return
	} else {
		//错误处理 只不过提前在web端控制台设置好就不会出现这个问题

	}

	//4、publish
	err = rabchannel.Publish(
		config.MqUploadExchangeName,
		routekey,
		true,  //如果为 ture 如果无法找到符合条件的队列 那么返回信息给发送者
		false, //不起作用？
		amqp.Publishing{
			ContentType: "text/plain", //明文格式
			Body:        []byte(msg),  //发送的信息 只能是字节信息
			Expiration:  "5000",       //ttl 这是5s
		})
	if err != nil {
		//错误处理 这个可以是转到死信交换机
		return
	}
	return

}

func Redispublish(routekey, msg string) {
	//如果在web控制台创建了exchange 1可以不用  发布端只负责发布到交换机
	//1、先定义出exchange
	err := rabchannel.ExchangeDeclare( //??? 属性细看
		"filestore-redis",
		"direct", //路由类型
		true,     //是否持久化
		false,    //是否丢失自动删除
		false,    //是否具有排他性
		false,    //是否堵塞
		nil,      //额外属性
	)
	if err != nil {
		fmt.Println("rabconn channel err:", err)
		return
	} else {
		//错误处理 只不过提前在web端控制台设置好就不会出现这个问题

	}

	//4、publish
	err = rabchannel.Publish(
		"filestore-redis",
		routekey,
		true,  //如果为 ture 如果无法找到符合条件的队列 那么返回信息给发送者
		false, //不起作用？
		amqp.Publishing{
			ContentType: "text/plain", //明文格式
			Body:        []byte(msg),  //发送的信息 只能是字节信息
			Expiration:  "5000",       //ttl 这是5s
		})
	if err != nil {
		//错误处理 这个可以是转到死信交换机
		return
	}
	return

}
