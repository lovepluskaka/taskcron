package taskcron

import (
	"strconv"
	"time"
)

const (
	Task_Status_Not_Executing = 1 // 未执行
	Task_Status_Executing     = 2 // 执行中
	Task_Status_Executed      = 3 // 已执行
	Task_Status_Cancel        = 4 // 已取消
	Task_Status_Fail          = 5 // 执行任务失败
)

const (
	Task_Execute_Method_Get    = 1 // Get请求
	Task_Execute_Method_Post   = 2 // Post请求
	Task_Execute_Method_Put    = 3 // Put请求
	Task_Execute_Method_Delete = 4 // Delete请求
)

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

func (t *TaskModel) toMap() map[string]interface{} {
	result := make(map[string]interface{})
	result["id"] = t.Id
	result["createTime"] = t.CreateTime.String()
	result["updateTime"] = t.UpdateTime.String()
	result["nextTime"] = t.NextTime.String()
	result["status"] = t.Status
	result["method"] = t.Method
	result["url"] = t.Url
	return result
}

func mapToStruct(m map[string]string) *TaskModel {
	id, _ := strconv.ParseInt(m["id"], 10, 64)
	ct, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", m["createTime"])
	ut, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", m["updateTime"])
	nt, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", m["nextTime"])
	status, _ := strconv.Atoi(m["status"])
	method, _ := strconv.Atoi(m["method"])

	result := &TaskModel{
		Id:         uint64(id),
		CreateTime: ct,
		UpdateTime: ut,
		NextTime:   nt,
		Status:     status,
		Method:     method,
		Url:        m["url"],
	}

	return result
}
