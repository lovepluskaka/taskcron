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
		Status:     Task_Status_Not_Executing,
		Method:     method,
		Url:        url,
	}

	// 开启事务，创建
	tx := t.redis.TxPipeline()

	if err := tx.HMSet(t.getTaskKey(id), task.toMap()).Err(); err != nil {
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
	lockKey := t.GetLockKey(id)

	// 执行任务前判断是否存在锁
	success, err := t.redis.SetNX(lockKey, 1, time.Duration(t.lockExpire*1000*1000)).Result()
	if err != nil {
		log.Fatal("任务加锁失败：", lockKey)
	}

	// 若锁已存在/加锁失败，1000ms后重试
	if !success {
		// 500ms后重试
		time.Sleep(1000 * 1000 * 1000)
		return t.do(id)
	} else {
		return t.doing(id)
	}
}

// 任务正在执行
func (t *Task) doing(id int64) error {

	taskKey := t.getTaskKey(id)

	taskMap, err := t.redis.HGetAll(taskKey).Result()

	if err != nil {
		return err
	}

	task := mapToStruct(taskMap)

	// 任务的状态必须是未执行
	if task.Status != Task_Status_Not_Executing {
		return TaskHasExcuteError
	}

	// 将任务状态改为执行中
	if err := t.redis.HSet(taskKey, "status", Task_Status_Executing).Err(); err != nil {
		return err
	}

	// 执行任务
	executeRes := get(task.Url)

	// 执行任务完成之后的操作
	if err := t.done(id, taskKey, executeRes); err != nil {
		return err
	}

	return nil
}

// 定时任务执行完成之后的操作
func (t *Task) done(id int64, key string, doneRes bool) error {

	data := make(map[string]interface{})

	data["updateTime"] = time.Now().String()
	if doneRes {
		data["status"] = Task_Status_Executed
	} else {
		data["status"] = Task_Status_Fail
	}

	// 开启事务，执行完成
	tx := t.redis.TxPipeline()

	// 更新任务状态
	if err := tx.HMSet(key, data).Err(); err != nil {
		tx.Close()
		return err
	}

	// 将任务从未执行的任务中移除
	tx.SRem(t.mapKey, id)

	// 删除悲观锁
	tx.Del(t.GetLockKey(id))

	tx.Exec()

	return nil
}

// 取消任务
func (t *Task) Cancel() error {
	return nil
}

// 获取任务key
func (t *Task) getTaskKey(id int64) string {
	return fmt.Sprintf("tasks:%s:%d", t.prefix, id)
}

// 获取锁的Key
func (t *Task) GetLockKey(id int64) string {
	return fmt.Sprintf("tasks:lock:%d", id)
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
