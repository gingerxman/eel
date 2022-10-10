package event

import (
	"context"
	"fmt"
	"github.com/gingerxman/eel/config"
	"github.com/gingerxman/eel/log"
	"github.com/gingerxman/gorm"
	"runtime/debug"
)

// IRestContext 解决循环依赖的临时方案
type IRestContext interface {
	AddEvent(eventData map[string]interface{})
}

type localEngine struct{}

func (this *localEngine) GetType() string {
	return "local"
}

func (this *localEngine) Send(ctx context.Context, eventData map[string]interface{}, tag string) {
	rc := ctx.Value("restContext")
	if rc == nil {
		return
	}
	restContext := rc.(IRestContext)

	restContext.AddEvent(eventData)
}

func newLocalEngine() *localEngine {
	return new(localEngine)
}

func EmitLocalEvent(eventData map[string]interface{}) {
	eventName := eventData["_event_name"].(string)
	for _, h := range getHandlersForEvent(eventName) {
		log.Logger.Info(fmt.Sprintf("[event] %s emit", eventName))
		go func(handler localEventHandler) {
			defer func() {
				if err := recover(); err != nil {
					log.Logger.Error(string(debug.Stack()))
				}
			}()

			ctx := newHandlerCtx()
			err := handler.Handle(ctx, eventData["data"].(map[string]interface{}))
			if err != nil {
				log.Logger.Error(err)
			}
		}(h)
	}
}

func newHandlerCtx() context.Context {
	ctx := context.Background()
	enableDb := config.ServiceConfig.DefaultBool("db::ENABLE_DB", true)
	var o *gorm.DB
	if enableDb {
		o = config.Runtime.DB
		ctx = context.WithValue(ctx, "orm", o)
	}

	return ctx
}

// localEngine发出事件的处理器
type localEventHandler interface {
	GetEventName() string
	Handle(ctx context.Context, eventData map[string]interface{}) error
}

var event2Handlers = make(map[string][]localEventHandler)

func getHandlersForEvent(eventName string) []localEventHandler {
	return event2Handlers[eventName]
}

func RegisterEventHandler(handler localEventHandler) {
	if v, ok := event2Handlers[handler.GetEventName()]; ok {
		event2Handlers[handler.GetEventName()] = append(v, handler)
	} else {
		event2Handlers[handler.GetEventName()] = []localEventHandler{handler}
	}
}

func init() {
	if config.ServiceConfig.DefaultBool("event::ENABLE_LOCAL_ENGINE", false) {
		registerEngine(newLocalEngine())
	}
}
