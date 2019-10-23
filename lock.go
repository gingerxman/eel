package eel

import (
	"fmt"
	"github.com/gingerxman/eel/config"
	"github.com/gingerxman/eel/log"
	"github.com/gomodule/redigo/redis"
	"errors"
	"time"
	"github.com/go-redsync/redsync"
)

type ILock interface {
	Lock(key string) (*redsync.Mutex, error)
}

//DummyLock 空的锁引擎
type DummyLock struct {

}

func (this *DummyLock) Lock(key string) (*redsync.Mutex, error) {
	log.Logger.Debug(fmt.Sprintf("[lock] lock by dummy engine : %s", key))
	return nil, nil
}

//RedisLock 基于redis的锁引擎
type RedisLock struct {
	engine *redsync.Redsync
}

func (this *RedisLock) Lock(key string) (*redsync.Mutex, error) {
	log.Logger.Debug(fmt.Sprintf("[lock] lock by redis engine : %s", key))
	if this.engine == nil {
		log.Logger.Warn("[lock] redsync engine is nil")
		return nil, nil
	} else {
		mutex := this.engine.NewMutex(key, redsync.SetExpiry(10*time.Second))
		err := mutex.Lock()
		if err != nil {
			log.Logger.Error(err)
			return nil, err
		} else {
			return mutex, nil
		}
	}
}


var Lock ILock //暴露的锁

var lockRedisAddress string = ""
var lockDbNum int = 1
var lockRedisPassword string = ""
var lockRedisPool *redis.Pool = nil


func init() {
	lockEngine := config.ServiceConfig.String("lock::ENGINE")
	if lockEngine == "" {
		lockEngine = "dummy"
	}
	
	if lockEngine == "dummy" {
		log.Logger.Info("[lock] use DummyLock")
		Lock = new(DummyLock)
	} else {
		lockRedisAddress = config.ServiceConfig.String("lock::REDIS_ADDRESS")
		lockDbNum, _ = config.ServiceConfig.Int("lock::REDIS_DB")
		lockRedisPassword = config.ServiceConfig.String("lock::REDIS_PASSWORD")
		log.Logger.Info(fmt.Sprintf("[lock] use RedisLock: %s - %d", lockRedisAddress, lockDbNum))
		
		// initialize a new pool
		lockRedisPool = &redis.Pool{
			MaxIdle:     10,
			IdleTimeout: 180 * time.Second,
			Dial: func() (c redis.Conn, err error) {
				if lockRedisAddress == "" {
					return nil, errors.New("invalid redisAddress")
				}
				
				c, err = redis.Dial("tcp", lockRedisAddress)
				if err != nil {
					log.Logger.Error(err)
					return nil, err
				}
				
				if lockRedisPassword != "" {
					if _, err := c.Do("AUTH", lockRedisPassword); err != nil {
						log.Logger.Error(err)
						c.Close()
						return nil, err
					}
				}
				
				_, selecterr := c.Do("SELECT", lockDbNum)
				if selecterr != nil {
					log.Logger.Error(selecterr)
					c.Close()
					return nil, selecterr
				}
				return
			},
			MaxConnLifetime: 60 * time.Minute,
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}
				_, err := c.Do("PING")
				return err
			},
		}
		
		//pool热身
		c := lockRedisPool.Get()
		defer c.Close()
		
		//创建
		Lock = &RedisLock{
			engine: redsync.New([]redsync.Pool{lockRedisPool}),
		}
	}
}