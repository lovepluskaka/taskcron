package taskcron

import (
	"github.com/go-redis/redis"
)

// redis 相关的配置
type RedisOptions *redis.Options

// task 任务相关的配置
type TaskOptions struct {
	Prefix     string // 任务在redis中存储的前缀
	LockExpire uint32 // 悲观（互斥）锁的过期时间
	Expire     uint32 // 任务执行完成后在redis中存储的过期时间
}

// 任务客户端
type Task struct {
	redis      *redis.Client
	prefix     string
	lockExpire uint32
	expire     uint32
}

func (t *Task) Create() {

}

// 初始化任务客户端
func Init(option RedisOptions, taskOptions TaskOptions) *Task {
	client := redis.NewClient(option)

	return &Task{
		redis:      client,
		prefix:     taskOptions.Prefix,
		lockExpire: taskOptions.LockExpire,
		expire:     taskOptions.Expire,
	}
}
