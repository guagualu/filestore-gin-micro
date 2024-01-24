package data

import (
	"context"
	"fileStore/conf"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"sync"
	"time"
)

type Data struct {
	db  *gorm.DB
	red *redis.Pool
}

// singel instance
var data Data
var once sync.Once

func GetData() *Data {

	once.Do(func() {
		data = Data{}
		data.db = NewDB(conf.GetConfig())
		data.red = NewRedis(conf.GetConfig())
	})
	return &data
}

// 用来承载事务的上下文
type contextTxKey struct{}

func NewDB(conf conf.Conf) *gorm.DB {
	// 终端打印输入 sql 执行记录
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  logger.Error,
		},
	)
	db, err := gorm.Open(mysql.Open(conf.DbConfig.Resource), &gorm.Config{
		DisableNestedTransaction:                 true,
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名不加s
		},
	})

	if err != nil {
		panic("failed to connect database")
	}
	if err = db.AutoMigrate(
		//User{},
		//File{},
		//UserFile{},
		//imSession{},
		//imSessionContent{},
		Friends{},
	); nil != err {
		panic("failed auto migrate")
	}

	return db
}

func NewRedis(conf conf.Conf) *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle:     10,
		MaxActive:   20,
		IdleTimeout: conf.RedisConfig.DialTimeout,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.RedisConfig.Addr, redis.DialPassword(conf.RedisConfig.Password), redis.DialDatabase(int(conf.RedisConfig.Db)), redis.DialReadTimeout(time.Second*5), redis.DialConnectTimeout(5*time.Second))
			if nil != err {
				return nil, err
			}
			return c, nil
		},
	}
	conn := redisPool.Get()
	if _, err := conn.Do("PING"); nil != err {
		panic(err)
	}

	return redisPool
}

// ExecTx gorm Transaction
func (d *Data) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

// DB 根据此方法来判断当前的db是不是使用事务的DB
func (d *Data) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.db
}

func (d *Data) RDB() *redis.Pool {
	return d.red
}

// 获取redis分布式锁
const (
	lockCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
    return "OK"
else
    return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
end`
	delCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`
)
const redisLock = "redisLock"

type WatchRes struct {
	Res bool
	Err error
}

func SetMutex(uuid string, ctx context.Context) error {
	//key count 是要输入的参数中key的数量
	db := GetData()
	lua := redis.NewScript(1, lockCommand)
	conn, err := db.red.Dial()
	if err != nil {
		return err
	}
	defer conn.Close()
	// uuid 以及 超时时间
	res, err := redis.String(lua.Do(conn, redisLock, uuid, 1500))
	if err != nil {
		return err
	}
	if res == "OK" {
		go db.watchAndPX(ctx, conn, uuid)
	}
	return nil
}

func (db *Data) watchAndPX(ctx context.Context, conn redis.Conn, uuid string) {
	//使用定时器 进行续约
	setExTimer := time.NewTimer(1000)
	defer setExTimer.Stop()
	for {
		select {
		case <-setExTimer.C:
			fmt.Println("------------------------------------")
			redis.Int(conn.Do("SET", redisLock, uuid, "PX", "1000"))
		case <-ctx.Done():
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			return
		}
	}
}

func DeleteMutex(uuid string) (bool, error) {
	//key count 是要输入的参数中key的数量
	db := GetData()
	lua := redis.NewScript(1, delCommand)
	conn, err := db.red.Dial()
	if err != nil {
		return false, err
	}
	defer conn.Close()
	// uuid 以及 超时时间
	res, err := redis.Int(lua.Do(conn, redisLock, uuid))
	if err != nil {
		return false, err
	}
	if res == 1 {
		return true, nil
	}
	return false, nil
}
