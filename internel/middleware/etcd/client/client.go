package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

const schema = "etcd"

func NewConnection(endPoints string, timeOut time.Duration, service string) (*grpc.ClientConn, error) {
	r := NewResolver(
		strings.Split(endPoints, ";"), service)
	resolver.Register(r)
	//ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	//defer cancel()

	addr := fmt.Sprintf("%s:///%s", r.Scheme(), service)
	// addr := service
	fmt.Println(addr)
	conn, err := grpc.DialContext(context.Background(), addr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		//指定初始化round_robin => balancer (后续可以自行定制balancer和 register、resolver 同样的方式)
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		grpc.WithBlock())
	if err != nil {
		fmt.Println("err:", err)
		return nil, err
	}
	fmt.Println("2")
	return conn, nil
}

// Resolver 实现grpc的grpc.resolve.Builder接口的Build与Scheme方法
type Resolver struct {
	// endpoints Etcd分布式存储结点地址
	endpoints []string
	// service 服务名称
	service string
	// cli Etcd客户端提供并管理一个 etcd v3 客户端会话。
	cli *clientv3.Client
	// cc ClientConn 包含解析器的回调，用于通知 gRPC ClientConn 的任何更新
	cc resolver.ClientConn
}

// NewResolver 返回一个resolver.Builder对象，endpoints为etcd服务器地址，service为服务名
func NewResolver(endpoints []string, service string) resolver.Builder {
	return &Resolver{endpoints: endpoints, service: service}
}

// Scheme 返回自定义的etcd服务格式化字符串
func (r *Resolver) Scheme() string {
	return schema
}

// ResolveNow 实现grpc.resolve.Builder接口的方法  实现包含这个方法的所有方法 结构体就实现这个接口变为buider
func (r *Resolver) ResolveNow(rn resolver.ResolveNowOptions) {
}

// Close 实现grpc.resolve.Builder接口的方法，用来关闭etcd客户端连接
func (r *Resolver) Close() {
	r.cli.Close()
}

// Build 实现grpc.resolve.Builder接口的方法，当调用`grpc.Dial()`时执行
func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   r.endpoints,
		DialTimeout: 5 * time.Second,
		Username:    "root",
		Password:    "meixi253",
	})
	if err != nil {
		return nil, fmt.Errorf("grpclb: create clientv3 client failed: %v", err)
	}
	r.cli = cli
	r.cc = cc
	prefix := "/" + target.URL.Scheme + "/" + r.service + "/"
	//用于同步 等待初次的reslover 更新完成
	sign := make(chan int)
	go r.watch(prefix, sign) //注册中心有变化 就更新在服务端的  获取 、更新 、赋值在map中 都集成在watch中了
	<-sign
	return r, nil
}

// watch 监听服务改动并更新服务列表
func (r *Resolver) watch(prefix string, sign chan int) {
	addrDict := make(map[string]resolver.Address)
	// 更新服务列表
	update := func() {
		addrList := make([]resolver.Address, 0, len(addrDict))
		for _, v := range addrDict {
			addrList = append(addrList, v)
		}
		r.cc.UpdateState(resolver.State{ //更新加载进resolver.ClientConn里
			Addresses: addrList,
		})
	}

	// 获取服务列表
	resp, err := r.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err == nil {
		for _, kv := range resp.Kvs {
			addrDict[string(kv.Key)] = resolver.Address{Addr: string(kv.Value)}
		}
	}
	update()
	//初始化完成 可以返回
	sign <- 1
	for {
		// 监听更新服务列表
		rch := r.cli.Watch(context.Background(), prefix, clientv3.WithPrefix(), clientv3.WithPrevKV()) //当监听的值变化时 就会返回这个channel

		for n := range rch { //循环阻塞监听
			for _, ev := range n.Events {
				switch ev.Type {
				case mvccpb.PUT:
					addrDict[string(ev.Kv.Key)] = resolver.Address{Addr: string(ev.Kv.Value)}
				case mvccpb.DELETE:
					delete(addrDict, string(ev.PrevKv.Key))
				}
			}
			update()
		}
	}
}
