package taskcron

import "time"

const (
	Task_Status_Not_Executing = 1 // 未执行
	Task_Status_Executing     = 2 // 执行中
	Task_Status_Executed      = 3 // 已执行
	Task_Status_Cancel        = 4 // 已取消
)

// 任务模型
type taskModel struct {
	Id         int32     // 任务id
	Version    int32     // 版本锁
	CreateTime time.Time // 创建时间
	UpdateTime time.Time // 更新时间
	LastTime   time.Time // 上次执行时间
	NextTime   time.Time // 下次执行时间
	Status     int       // 任务执行状态
}

func (t *taskModel) toMap() map[string]interface{} {
	result := make(map[string]interface{})
	result["id"] = t.Id
	result["version"] = t.Version
	result["createTime"] = t.CreateTime.String()
	result["updateTime"] = t.UpdateTime.String()
	result["lastTime"] = t.LastTime.String()
	result["nextTime"] = t.NextTime.String()
	return result
}
