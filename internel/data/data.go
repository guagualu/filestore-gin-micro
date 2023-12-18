package data

import (
	"context"
	"fileStore/conf"
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

func GetData() Data {

	once.Do(func() {
		data = Data{}
		data.db = NewDB(conf.GetConfig())
		data.red = NewRedis(conf.GetConfig())
	})
	return data
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
	if err = db.AutoMigrate(); nil != err {
		panic("failed auto migrate")
	}

	return db
}

func NewRedis(conf conf.Conf) *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle:     10,
		MaxActive:   20,
		IdleTimeout: conf.RedisConfig.DialTimeout,
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
