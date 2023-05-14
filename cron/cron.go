package cron

import (
	"context"
	"fmt"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/eel/config"
	"github.com/gingerxman/eel/cron/toolbox"
	"github.com/gingerxman/eel/handler"
	"github.com/gingerxman/eel/log"
	"github.com/gingerxman/gorm"
	"runtime/debug"
	"time"
)

type CronTask struct {
	name            string
	spec            string
	taskFunc        toolbox.TaskFunc
	onlyRunThisTask bool
}

func (this *CronTask) OnlyRun() {
	this.onlyRunThisTask = true
}

var name2task = make(map[string]*CronTask)

func newTaskCtx() *TaskContext {
	inst := new(TaskContext)
	ctx := context.Background()
	enableDb := config.ServiceConfig.DefaultBool("db::ENABLE_DB", true)
	var o *gorm.DB
	if enableDb {
		o = config.Runtime.DB
		ctx = context.WithValue(ctx, "orm", o)
	}

	resource := GetManagerResource(ctx)
	ctx = context.WithValue(ctx, "jwt", resource.CustomJWTToken)
	userId, authUserId, _ := eel.ParseUserIdFromJwtToken(resource.CustomJWTToken)
	ctx = context.WithValue(ctx, "user_id", userId)
	ctx = context.WithValue(ctx, "uid", authUserId)
	resource.Ctx = ctx
	inst.Init(ctx, o, resource)
	return inst
}

func taskWrapper(task taskInterface) toolbox.TaskFunc {

	return func() error {
		taskCtx := newTaskCtx()
		o := taskCtx.GetOrm()
		ctx := taskCtx.GetCtx()

		defer handler.RecoverFromCronTaskPanic(ctx)
		var fnErr error
		taskName := task.GetName()
		startTime := time.Now()
		log.Logger.Info(fmt.Sprintf("[%s] run...", taskName))
		if o != nil && task.IsEnableTx() {
			o.Begin()
			fnErr = task.Run(taskCtx)
			o.Commit()
		} else {
			fnErr = task.Run(taskCtx)
		}
		dur := time.Since(startTime)
		log.Logger.Info(fmt.Sprintf("[%s] done, cost %g s", taskName, dur.Seconds()))
		return fnErr
	}
}

func fetchData(pi pipeInterface) {
	taskName := pi.(taskInterface).GetName()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Logger.Warn(string(debug.Stack()))
				fetchData(pi)
			}
		}()
		for {
			data := pi.GetData()
			if data != nil {
				taskCtx := newTaskCtx()
				log.Logger.Info(fmt.Sprintf("[%s] consume data...", taskName))
				startTime := time.Now()
				pi.RunConsumer(data, taskCtx)
				dur := time.Since(startTime)
				log.Logger.Info(fmt.Sprintf("[%s] consume done, cost %g s !", taskName, dur.Seconds()))
			}
		}
	}()
}

func RegisterPipeTask(pi pipeInterface, spec string) *CronTask {
	task := RegisterTask(pi.(taskInterface), spec)
	if task != nil {
		if pi.EnableParallel() { // 并行模式下，开启通道容量十分之一的goroutine消费通道
			for i := pi.GetConsumerCount(); i > 0; i-- {
				fetchData(pi)
			}
		} else {
			fetchData(pi)
		}
	}
	return task
}

func RegisterTask(task taskInterface, spec string) *CronTask {
	if config.ServiceConfig.DefaultBool("system::ENABLE_CRON_MODE", false) || config.ServiceConfig.String("system::SERVICE_MODE") == "cron" {
		tname := task.GetName()
		wrappedFn := taskWrapper(task)
		cronTask := &CronTask{
			name:            tname,
			spec:            spec,
			taskFunc:        wrappedFn,
			onlyRunThisTask: false,
		}
		name2task[tname] = cronTask

		return cronTask
	} else {
		return nil
	}
}

func RegisterTaskInRestMode(task taskInterface, spec string) *CronTask {
	if !config.ServiceConfig.DefaultBool("system::ENABLE_CRON_MODE", false) && config.ServiceConfig.String("system::SERVICE_MODE") == "rest" {
		tname := task.GetName()
		wrappedFn := taskWrapper(task)
		cronTask := &CronTask{
			name:            tname,
			spec:            spec,
			taskFunc:        wrappedFn,
			onlyRunThisTask: false,
		}
		name2task[tname] = cronTask

		return cronTask
	} else {
		return nil
	}
}

func RegisterCronTask(tname string, spec string, f toolbox.TaskFunc) *CronTask {
	cronTask := &CronTask{
		name:            tname,
		spec:            spec,
		taskFunc:        f,
		onlyRunThisTask: false,
	}
	name2task[tname] = cronTask

	return cronTask
}

func StartCronTasks() {
	var onlyRunTask *CronTask
	for _, task := range name2task {
		if task.onlyRunThisTask {
			onlyRunTask = task
		}
	}

	if onlyRunTask != nil {
		cronTask := onlyRunTask
		log.Logger.Info("[cron] create cron task ", cronTask.name, cronTask.spec)
		task := toolbox.NewTask(cronTask.name, cronTask.spec, cronTask.taskFunc)
		toolbox.AddTask(cronTask.name, task)
	} else {
		for _, cronTask := range name2task {
			log.Logger.Info("[cron] create cron task ", cronTask.name, cronTask.spec)
			task := toolbox.NewTask(cronTask.name, cronTask.spec, cronTask.taskFunc)
			toolbox.AddTask(cronTask.name, task)
		}
	}

	toolbox.StartTask()
}

func StopCronTasks() {
	toolbox.StopTask()
}
