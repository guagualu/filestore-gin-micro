package server

import (
	"context"
	"fmt"
	"log"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const schema = "etcd"

// ServiceRegister 创建租约注册服务
type ServiceRegister struct {
	// cli etcd客户端
	cli *clientv3.Client
	// leaseID 租约ID
	leaseID clientv3.LeaseID
	// keepAliveChan 租约keepAlive相应chan
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	// key 放入etcd中对应的key
	key string
	// val 放入etcd中对应的value
	val string
}

// RegisterConfig 注册客户端配置
type RegisterConfig struct {
	// 继承etcd客户端配置
	clientv3.Config
	// ServerName 服务名称
	ServerName string
	// Address 服务地址
	Address string
	// Lease 租约时间TTL，以秒为单位的建议生存时间
	Lease int64
}

func init() {

}

// NewServiceRegister 新建注册服务
func NewServiceRegister(conf RegisterConfig) (*ServiceRegister, error) {
	cli, err := clientv3.New(conf.Config)
	// logger := pkg.NewLogger(pkg.LogConfig{SavePath: "./", FileName: "log", FileExt: ".txt", MaxSize: 1, MaxBackUps: 1, MaxAge: 1})
	if err != nil {
		log.Fatal("clinew,err:", err)
	}

	service := &ServiceRegister{
		cli: cli,
		key: "/" + schema + "/" + conf.ServerName + "/" + conf.Address, //规则 这么设 再按客户端这么找 才赵得到
		val: conf.Address,
		// logger: *logger,
	}
	//申请租约设置时间keepalive
	if err := service.putKeyWithLease(conf.Lease); err != nil {
		return nil, err
	}
	return service, nil
}

// putKeyWithLease 设置租约并放置key-value
func (s *ServiceRegister) putKeyWithLease(lease int64) error {
	//设置租约时间
	resp, err := s.cli.Grant(context.Background(), lease)
	if err != nil {
		return err
	}
	//注册服务并绑定租约
	_, err = s.cli.Put(context.Background(), s.key, s.val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}

	//设置续租 定期发送需求请求
	leaseRespChan, err := s.cli.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return err
	}

	s.leaseID = resp.ID
	s.keepAliveChan = leaseRespChan
	log.Printf("Put key:%s  val:%s  success!", s.key, s.val)
	return nil

}

// ListenLeaseRespChan 监听续租情况
func (s *ServiceRegister) ListenLeaseRespChan() {
	// for leaseKeepResp := range s.keepAliveChan {
	// 	// s.logger.Fatal("续约成功")
	// 	log.Fatal("续约成功", leaseKeepResp)
	// }
	// log.Println("关闭续租")
	for {
		select {
		case keepResp := <-s.keepAliveChan:
			if keepResp == nil {
				fmt.Println("租约已经失效了")
			} else { // 每秒会续租一次, 所以就会受到一次应答
				// s.logger.Fatal("续约成功")
				fmt.Println("收到自动续租应答:", keepResp.ID)
			}
		}
	}
}

// Close 注销服务
func (s *ServiceRegister) Close() error {
	//撤销租约
	if _, err := s.cli.Revoke(context.Background(), s.leaseID); err != nil {
		return err
	}
	log.Println("撤销租约")
	return s.cli.Close()
}
