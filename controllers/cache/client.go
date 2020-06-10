package cache

import (
	"fmt"
	rediscli "github.com/go-redis/redis"
	"strings"
	"time"
)

const (
	RedisDownScaling     = "RedisDownScaling"
	RedisUpScaling       = "RedisUpScaling"
	RedisRollingUpgrades = "RollingUpgrades"

	MessageRedisCluster = "Redis  %s already created."

	UpdateMessageRedisCluster = "Redis  %s already update."

	MessageRedisDownScaling     = "Redis downscale from %d to %d"
	MessageRedisUpScaling       = "Redis upscale from %d to %d"
	MessageRedisRollingUpgrades = "Redis resource from %s to %s"

	RedisSentinelConnPort  = "26379"
	RedisSentinelConnGroup = "mymaster"
)

type RedisConnect struct {
	Endpoint  string
	Port      string
	Password  string
	GroupName string
}

// NewRedisConnection returns redis connection
func NewRedisConnection(endpoint, port, password, groupName string) *RedisConnect {
	return &RedisConnect{
		Endpoint:  endpoint,
		Port:      port,
		Password:  password,
		GroupName: groupName,
	}
}

// NewRedisPool returns redis sentinel client
func (c *RedisConnect) NewRedisPool() *rediscli.Client {

	return BuildRedisPool(c.Endpoint, c.Port, c.Password, c.GroupName, 0)
}

// NewRedisClient returns redis client
func (c *RedisConnect) NewRedisClient() *rediscli.Client {

	return BuildRedisClient(c.Endpoint, c.Port, c.Password, 0)
}

// BuildRedisPool returns redis connection pool client
func BuildRedisPool(redisSentinelIP, redisSentinelPort, redisSentinelPassword, redisGroupName string, redisIndex int) *rediscli.Client {

	var sentinelsInfo []string
	sentinels := strings.Split(redisSentinelIP, ",")
	if len(sentinels) > 0 {
		for _, s := range sentinels {
			sp := s + ":" + redisSentinelPort
			sentinelsInfo = append(sentinelsInfo, sp)
		}
	}

	options := &rediscli.FailoverOptions{
		MasterName:         redisGroupName,
		SentinelAddrs:      sentinelsInfo,
		Password:           redisSentinelPassword,
		DB:                 redisIndex,
		PoolSize:           100,
		DialTimeout:        10 * time.Second,
		ReadTimeout:        30 * time.Second,
		WriteTimeout:       30 * time.Second,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        time.Millisecond,
		IdleCheckFrequency: time.Millisecond,
	}

	client := rediscli.NewFailoverClient(options)

	return client

}

// BuildRedisClient returns redis connection client
func BuildRedisClient(host, port, password string, index int) *rediscli.Client {

	options := &rediscli.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       index,
	}
	client := rediscli.NewClient(options)

	return client

}
