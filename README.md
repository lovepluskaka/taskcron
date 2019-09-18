# taskcron
基于redis的用于分布式系统的定时任务管理模块

### 使用方式

##### go mod
```$xslt
go mod download https://github.com/lovepluskaka/taskcron
```

##### go get
```$xslt
go get -u https://github.com/lovepluskaka/taskcron
```


### Example

```$xslt
    import (
      task https://github.com/lovepluskaka/taskcron
    )
    
    // 初始化任务管理
    t := task.Init(task.redisOptions, tassk.taskOptions)
    
    // 创建任务
    taskModel, err := t.Create(durtion, "http://127.0.0.1:3001/test", task.Task_Execute_Method_Get)
    
    // 删除任务
    
    err := t.Cancel(taskId)
```

### Tips

##### task.RedisOptions (redis.Options)

##### task.TaskOptions

```$xslt
// task 任务相关的配置
type TaskOptions struct {
	Prefix     string // 任务在redis中存储的前缀
	LockExpire uint32 // 悲观（互斥）锁的过期时间
	Expire     uint32 // 任务执行完成后在redis中存储的过期时间
}
```

##### task.TaskModel

```$xslt
// 任务模型
type TaskModel struct {
	Id         uint64    // 任务id
	CreateTime time.Time // 创建时间
	UpdateTime time.Time // 更新时间
	NextTime   time.Time // 下次执行时间
	Status     int       // 任务执行状态
	Method     int       // 任务的执行方法
	Url        string    // 请求的url地址
}
```