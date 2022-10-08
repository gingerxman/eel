package event

import (
	"context"
	"fmt"
	"github.com/gingerxman/eel/config"
)

type consoleEngine struct{}

func (this *consoleEngine) GetType() string {
	return "console"
}

func (this *consoleEngine) Send(ctx context.Context, data map[string]interface{}, tag string) {
	eventName := data["_event_name"]
	fmt.Printf("[Event] CONSOLE ENGINE: receive event %s with tag: %s", eventName, tag)
}

func newConsoleEngine() *consoleEngine {
	return new(consoleEngine)
}

func init() {
	if config.ServiceConfig.DefaultBool("event::ENABLE_CONSOLE_ENGINE", false) {
		registerEngine(newConsoleEngine())
	}
}
