package redis

import (
	"context"
	"fmt"
	"seabase/extend/conf"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

var redisConn *redis.Pool

// GetRedisConn 获取 Redis 客户端连接
func GetRedisConn() *redis.Pool {
	return redisConn
}

// Setup 创建 Redis 连接
func Init() error {
	redisConn = &redis.Pool{
		MaxIdle:     conf.RedisConf.MaxIdle,
		MaxActive:   conf.RedisConf.MaxActive,
		IdleTimeout: conf.RedisConf.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.RedisConf.Host+":"+strconv.Itoa(conf.RedisConf.Port))
			if err != nil {
				return nil, err
			}
			// 验证密码
			if conf.RedisConf.Password != "" {
				if _, err := c.Do("AUTH", conf.RedisConf.Password); err != nil {
					_ = c.Close()
					return nil, err
				}
			}
			// 选择数据库
			if _, err := c.Do("SELECT", conf.RedisConf.DBNum); err != nil {
				_ = c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	return nil
}

func Set(key string, data string, seconds int) error {
	conn := GetRedisConn().Get()
	defer conn.Close()
	
	_, err := conn.Do("SET", key, data)
	if err != nil {
		return err
	}
	
	_, err = conn.Do("EXPIRE", key, seconds)
	if err != nil {
		return err
	}
	return nil
}

func HSet(key string, field string, data string) error {
	conn := GetRedisConn().Get()
	defer conn.Close()
	
	_, err := conn.Do("HSET", key, field, data)
	if err != nil {
		return err
	}
	
	return nil
}

func HGet(key string, field string) (string, error) {
	conn := GetRedisConn().Get()
	defer conn.Close()
	
	reply, err := redis.String(conn.Do("HGET", key, field))
	if err != nil && err != redis.ErrNil {
		return "", err
	}
	if err == redis.ErrNil {
		return "", nil
	}
	return reply, nil
}

// 删除某字段数据
func HDel(key, field string) error {
	conn := GetRedisConn().Get()
	defer conn.Close()
	_, err := conn.Do("HDEL", key, field)
	return err
}

func HGetAll(key string) (map[string]string, error) {
	conn := GetRedisConn().Get()
	defer conn.Close()
	
	ret, err := redis.StringMap(conn.Do("HGETALL", key))
	if err != nil {
		return nil, err
	}
	
	return ret, nil
}

func Exists(key string) bool {
	conn := GetRedisConn().Get()
	defer conn.Close()
	
	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return exists
}

func Get(key string) (string, error) {
	conn := GetRedisConn().Get()
	defer conn.Close()
	
	reply, err := redis.String(conn.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		return "", err
	}
	if err == redis.ErrNil {
		return "", nil
	}
	return reply, nil
}

func GetTTL(key string) (int, error) {
	conn := GetRedisConn().Get()
	defer conn.Close()
	reply, err := redis.Int(conn.Do("TTL", key))
	if err != nil && err != redis.ErrNil {
		return 0, err
	}
	return reply, err
}

func Del(key string) (bool, error) {
	conn := GetRedisConn().Get()
	defer conn.Close()
	
	return redis.Bool(conn.Do("DEL", key))
}

func DelLike(key string) error {
	conn := GetRedisConn().Get()
	defer conn.Close()
	
	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}
	
	for _, key := range keys {
		_, err := Del(key)
		if err != nil {
			return err
		}
	}
	return nil
}

func ListenPubSubChannels(
	ctx context.Context, onStart func() error,
	onMessage func(channel string, data []byte) error, channels ...string) error {
	// A ping is set to the server with this period to test for the health of
	// the connection and server.
	const healthCheckPeriod = time.Minute * 5
	var err error
	conn := GetRedisConn().Get()
	defer conn.Close()
	
	var pscs []redis.PubSubConn
	psc := redis.PubSubConn{Conn: conn}
	pscs = append(pscs, psc)
	
	defer func() {
		for _, psc := range pscs {
			err = psc.Unsubscribe() // unsubscribe with no args unsubs all channels
			if err != nil {
				fmt.Println("----Unsubscribe failed----")
				fmt.Println(err)
			}
		}
	}()
	
	if err := psc.Subscribe(redis.Args{}.AddFlat(channels)...); err != nil {
		return err
	}
	
	done := make(chan error, 1)
	
	// Start a goroutine to receive notifications from the server.
	go func() {
		defer psc.Close()
		for {
			switch n := psc.Receive().(type) {
			case error:
				done <- n
				return
			case redis.Message:
				if err := onMessage(n.Channel, n.Data); err != nil {
					done <- err
					return
				}
			case redis.Subscription:
				switch n.Count {
				case len(channels):
					// Notify application when all channels are subscribed.
					if err := onStart(); err != nil {
						done <- err
						return
					}
				case 0:
					// Return from the goroutine when all channels are unsubscribed.
					done <- nil
					return
				}
			}
		}
	}()
	
	ticker := time.NewTicker(healthCheckPeriod)
	defer ticker.Stop()
	//loop:
	for {
		select {
		case <-ticker.C:
			// Send ping to test health of connection and server. If
			// corresponding pong is not received, then receive on the
			// connection will timeout and the receive goroutine will exit.
			if err = psc.Ping(""); err != nil {
				//break loop
				return err
			}
		case <-ctx.Done():
			if err := psc.Unsubscribe(); err != nil {
				return fmt.Errorf("redis pubsub unsubscribe err: %v", err)
			}
			//break loop
		case err := <-done:
			// Return error from the receive goroutine.
			return err
		}
		return nil
	}
	
	// Signal the receiving goroutine to exit by unsubscribing from all channels.
	//psc.Unsubscribe()
	
	// Wait for goroutine to complete.
	//return <-done
}

func UnSub(channel string) error {
	conn := GetRedisConn().Get()
	defer conn.Close()
	psc := redis.PubSubConn{Conn: conn}
	err := psc.Unsubscribe(channel)
	return err
}

func Publish(channel string, value string) error {
	conn := GetRedisConn().Get()
	defer conn.Close()
	
	_, err := conn.Do("PUBLISH", channel, value)
	if err != nil {
		return err
	}
	return nil
}

func LRange(key string, start int, end int) (interface{}, error) {
	conn := GetRedisConn().Get()
	defer conn.Close()
	return conn.Do("LRANGE", key, strconv.Itoa(start), strconv.Itoa(end))
}

func Keys(patten string) ([]string, error) {
	conn := GetRedisConn().Get()
	defer conn.Close()
	
	return redis.Strings(conn.Do("KEYS", patten))
}

func SMembers(key string) ([]interface{}, error) {
	conn := GetRedisConn().Get()
	defer conn.Close()
	return redis.Values(conn.Do("smembers", key))
}

func SRemoveMember(key string, member string) (interface{}, error) {
	conn := GetRedisConn().Get()
	defer conn.Close()
	return conn.Do("srem", member)
}
