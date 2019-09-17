package taskcron

import (
	"fmt"
	"github.com/go-redis/redis"
	"sync"
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
	redis      *redis.Client       // redis 客户端
	prefix     string              // 标示任务的前缀
	lockExpire uint32              // 锁
	expire     uint32              // 任务执行完成之后过期的时间
	idsKey     string              // 索引前缀
	mapKey     string              // 存储待执行任务的key
	timingMap  map[int]*time.Timer // 存储定时器的map
	sync.RWMutex
}

// 创建任务
func (t *Task) Create(d time.Duration, f func()) (*TaskModel, error) {
	id, err := t.redis.Incr(t.idsKey).Result()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	task := &TaskModel{
		Id:         uint64(id),
		Version:    1,
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
	time.AfterFunc(d, f)
	return task, nil
}

// 执行任务
func (t *Task) Do() (*TaskModel, error) {
	return nil, nil
}

// 完成任务
func (t *Task) Done() (*TaskModel, error) {
	return nil, nil
}

// 取消任务
func (t *Task) Cancel() error {
	return nil
}

// 添加定时器
func (t *Task) Set(k int, v *time.Timer) {
	t.Lock()
	t.timingMap[k] = v
	defer t.Unlock()
}

// 获取定时器
func (t *Task) Get(k int) *time.Timer {
	return t.timingMap[k]
}

// 删除定时器
func (t *Task) Delete(k int) {
	t.Lock()
	delete(t.timingMap, k)
	defer t.Unlock()
}

// 停止定时器
func (t *Task) Stop(k int) error {
	timming := t.Get(k)
	if timming != nil {
		timming.Stop()
		return nil
	} else {
		return TimmerNotFoundError
	}
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
		lockExpire: taskOptions.LockExpire,
		expire:     taskOptions.Expire,
		idsKey:     idKey,
		mapKey:     mapKey,
	}, nil
}
