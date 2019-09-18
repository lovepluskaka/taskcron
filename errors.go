package taskcron

import "errors"

var TaskHasExcuteError error = errors.New("任务不是未执行的状态")
