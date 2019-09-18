package taskcron

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"time"
)

// redis 相关的配置
type RedisOptions *redis.Options

// task 任务相关的配置
type TaskOptions struct {
	Prefix     string // 任务在redis中存储的前缀
	LockExpire uint32 // 悲观（互斥）锁的过期时间
	Expire     uint32 // 任务执行完成后在redis中存储的过期时间
}

// 返回结果
//type Result struct {
//	Id uint64 // 任务id
//	Version int32 // 版本锁
//	Status int // 任务状态
//}

// 任务客户端
type Task struct {
	redis      *redis.Client // redis 客户端
	prefix     string        // 标示任务的前缀
	lockExpire uint32        // 锁的过期时间
	expire     uint32        // 任务执行完成之后过期的时间
	idsKey     string        // 索引前缀
	mapKey     string        // 存储待执行任务的key
}

// 创建任务
func (t *Task) Create(d time.Duration, url string, method int) (*TaskModel, error) {
	id, err := t.redis.Incr(t.idsKey).Result()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	task := &TaskModel{
		Id:         uint64(id),
		CreateTime: now,
		UpdateTime: now,
		NextTime:   now.Add(d),
	}

	taskKey := fmt.Sprintf("tasks:%s:%d", t.prefix, id)

	// 开启事务，创建
	tx := t.redis.TxPipeline()

	if err := tx.HMSet(taskKey, task.toMap()).Err(); err != nil {
		tx.Close()
		return nil, err
	}
	if err := tx.SAdd(t.mapKey, id).Err(); err != nil {
		tx.Close()
		return nil, err
	}

	tx.Exec()

	// 将方法推送到定时任务中
	time.AfterFunc(d, func() {
		err := t.do(id)
		if err != nil {
			log.Fatalf("error is %s\n", err.Error())
		}
	})
	return task, nil
}

// 执行任务
func (t *Task) do(id int64) error {
	lockKey := fmt.Sprintf("tasks:lock:%d", id)

	// 执行任务前判断是否存在锁
	success, err := t.redis.SetNX(lockKey, 1, time.Duration(t.lockExpire*1000*1000)).Result()
	if err != nil {
		log.Fatal("任务加锁失败：", lockKey)
	}

	// 若锁已存在/加锁失败，500ms后重试
	if !success {
		// 500ms后重试
		time.Sleep(500 * 1000 * 1000)
		return t.do(id)
	} else {

		return nil
	}
}

// 完成任务
func (t *Task) Done() (*TaskModel, error) {

	return nil, nil
}

// 取消任务
func (t *Task) Cancel() error {
	return nil
}

// 初始化任务客户端
func Init(option RedisOptions, taskOptions TaskOptions) (*Task, error) {
	client := redis.NewClient(option)

	// 设置任务的id
	idKey := fmt.Sprintf("tasks:ids:%s", taskOptions.Prefix)

	if err := client.SetNX(idKey, 0, 0).Err(); err != nil {
		return nil, err
	}

	// 存储待执行任务的key
	mapKey := fmt.Sprintf("tasks:waiting:%s", taskOptions.Prefix)

	return &Task{
		redis:      client,
		prefix:     taskOptions.Prefix,
		lockExpire: taskOptions.LockExpire | 500,   // 默认锁时间500毫秒
		expire:     taskOptions.Expire | 3600*24*7, // 默认任务过期时间7天
		idsKey:     idKey,
		mapKey:     mapKey,
	}, nil
}
