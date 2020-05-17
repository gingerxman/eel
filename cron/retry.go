package cron

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gingerxman/eel/backoff"
	"github.com/gingerxman/eel/log"
	"runtime"
	"time"
)

type RetryTaskParam struct {
	NewContext func() context.Context
	GetDatas func() []interface{}
	BeforeAction func(data interface{}) error
	DoAction func(ctx context.Context, times int, data interface{}) error
	AfterActionSuccess func(data interface{}) error
	AfterActionFail func(data interface{}) error
	GetTaskDataId func(data interface{}) string
	RecordFailByPanic func(data interface{}, error string)
}

func retryTaskGorutione(maxMinutes int, taskParam *RetryTaskParam, goroutineTimes int, data interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Logger.Error(err)
			errMsg := fmt.Sprintf("%s", err)
			var buffer bytes.Buffer
			buffer.WriteString(fmt.Sprintf("[Unprocessed_Exception] %s\n", errMsg))
			for i := 1; ; i++ {
				_, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				}
				buffer.WriteString(fmt.Sprintf("%s:%d\n", file, line))
			}
			log.Logger.Error(buffer.String())
			
			if goroutineTimes <= 3 {
				time.Sleep(4 * time.Second)
				log.Logger.Warn(fmt.Sprintf("[retry] restart goroutine for %d times", goroutineTimes))
				go retryTaskGorutione(maxMinutes, taskParam, goroutineTimes+1, data)
			} else {
				//需要捕捉在AfterActionFail和
				defer func() {
					if err3 := recover(); err3 != nil {
						log.Logger.Error(err3)
					}
				}()
				
				err2 := taskParam.AfterActionFail(data)
				if err2 != nil {
					log.Logger.Error(err2)
				}
				
				taskParam.RecordFailByPanic(data, errMsg)
			}
		}
	}()
	
	expBackoff := &backoff.ExponentialBackOff{
		InitialInterval:     backoff.DefaultInitialInterval,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          1.8,// * backoff.DefaultMultiplier,
		MaxInterval:         backoff.DefaultMaxInterval,
		MaxElapsedTime:      time.Duration(maxMinutes) * time.Minute,
		Clock:               backoff.SystemClock,
	}
	expBackoff.Reset()
	
	if taskParam.BeforeAction != nil {
		err := taskParam.BeforeAction(data)
		if err != nil {
			log.Logger.Error(err)
			return
		}
	}
	
	times := 0
	ctx := context.Background()
	err := backoff.RetryNotify(func() error {
		times += 1
		return taskParam.DoAction(ctx, times, data)
	}, expBackoff, func (err error, duration time.Duration) {
		log.Logger.Warn(fmt.Sprintf("[push_order_payment] push '%s' fail %d times, because of : %s, next push after %v", taskParam.GetTaskDataId(data), times, err.Error(), duration))
	})
	
	if err != nil {
		log.Logger.Error(err)
		
		err := taskParam.AfterActionFail(data)
		if err != nil {
			log.Logger.Error(err)
		}
	} else {
		err := taskParam.AfterActionSuccess(data)
		if err != nil {
			log.Logger.Error(err)
		}
	}
}

func StartRetryTask(maxMinutes int, taskParam *RetryTaskParam) {
	if taskParam.BeforeAction == nil {
		log.Logger.Error("[retry] Need taskParam.BeforeAction != nil")
		return
	}
	if taskParam.DoAction == nil {
		log.Logger.Error("[retry] Need taskParam.DoAction != nil")
		return
	}
	if taskParam.GetDatas == nil {
		log.Logger.Error("[retry] Need taskParam.GetDatas != nil")
		return
	}
	if taskParam.AfterActionSuccess == nil {
		log.Logger.Error("[retry] Need taskParam.AfterActionSuccess != nil")
		return
	}
	if taskParam.AfterActionFail == nil {
		log.Logger.Error("[retry] Need taskParam.AfterActionFail != nil")
		return
	}
	if taskParam.GetTaskDataId == nil {
		log.Logger.Error("[retry] Need taskParam.GetTaskDataId != nil")
		return
	}
	if taskParam.RecordFailByPanic == nil {
		log.Logger.Error("[retry] Need taskParam.RecordFailByPanic != nil")
		return
	}
	
	//datas := getOrdersNeedPush()
	datas := taskParam.GetDatas()
	
	for _, data := range datas {
		go retryTaskGorutione(maxMinutes, taskParam, 1, data)
	}
}
