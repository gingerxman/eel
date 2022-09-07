package handler

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gingerxman/eel/config"
	"github.com/gingerxman/eel/log"
	"net/http"
	"runtime/debug"
	"strings"
)

var sentryChannel = make(chan map[string]interface{}, 2048)

func isEnableSentry() bool {
	return config.ServiceConfig.DefaultBool("sentry::ENABLE_SENTRY", false)
}

const SENTRY_CHANNEL_TIMEOUT = 50

// CapturePanicToSentry will collect error info then send to sentry
func CaptureErrorToSentry(ctx *Context, err string) {
	if !isEnableSentry() {
		log.Logger.Debug("[sentry] Sentry is not enabled, Please enable it!!!")
		//Logger.Infof("http server Running on http://%s\n", this.Server.Addr)
		//beegoMode := os.Getenv("BEEGO_RUNMODE")
		//if beegoMode == "prod" {
		//	Warn("Sentry is not enabled under prod mode, Please enable it!!!!")
		//}
		return
	}

	data := make(map[string]interface{})
	data["err_msg"] = err
	data["service_name"] = config.ServiceConfig.String("SERVICE_NAME")

	//skipFramesCount := AppConfig.DefaultInt("sentry::SKIP_FRAMES_COUNT", 3)
	//contextLineCount := AppConfig.DefaultInt("sentry::CONTEXT_LINE_COUNT", 5)
	//appRootPath := AppConfig.String("appname")
	//inAppPaths := []string{appRootPath}

	//var sStacktrace *raven.Stacktrace
	//var sError, ok = err.(error)
	//if ok {
	//	sStacktrace = raven.GetOrNewStacktrace(sError, skipFramesCount, contextLineCount, inAppPaths)
	//} else {
	//	sStacktrace = raven.NewStacktrace(skipFramesCount, contextLineCount, inAppPaths)
	//}
	//sStacktrace = raven.NewStacktrace(skipFramesCount, contextLineCount, inAppPaths)
	//sException := raven.NewException(sError, sStacktrace)
	//spew.Dump(sStacktrace)
	data["stack"] = string(debug.Stack())
	data["raven_http"] = raven.NewHttp(ctx.Request.HttpRequest)
	data["http_request"] = ctx.Request.HttpRequest

	select {
	case sentryChannel <- data:
	default:
		//metrics.GetSentryChannelTimeoutCounter().Inc()
		log.Logger.Warn("[sentry] push timeout")
		//case <-time.After(time.Millisecond * SENTRY_CHANNEL_TIMEOUT):
		//	metrics.GetSentryChannelTimeoutCounter().Inc()
		//	Warn("[sentry] push timeout")
	}

}

func sendSentryPacketV2(data map[string]interface{}) {
	var packet *raven.Packet
	errMsg := data["err_msg"].(string)

	//封装http request
	httpRequest, ok := data["http_request"].(*http.Request)
	if ok {
		ravenHttp := raven.NewHttp(httpRequest)

		method := strings.ToLower(httpRequest.Method)
		if method == "post" || method == "put" || method == "delete" {
			data := make(map[string]string)
			for key, _ := range httpRequest.PostForm {
				value := httpRequest.PostForm.Get(key)
				if len(value) >= 100 {
					value = value[:100] + "..."
				}
				data[key] = value
			}
			ravenHttp.Data = data
		}

		packet = raven.NewPacket(errMsg, ravenHttp)
	} else {
		packet = raven.NewPacket(errMsg)
	}

	//确定extra
	if extra, ok := data["extra"]; ok {
		packet.Extra = extra.(map[string]interface{})
	} else {
		packet.Extra = make(map[string]interface{})
	}

	//确定堆栈信息
	stack, ok := data["stack"].(string)
	if !ok {
		stack = "no stack"
	}
	packet.Extra["stacktrace"] = stack

	//其他Tag
	tags := map[string]string{
		"service_name": data["service_name"].(string),
	}

	//发送给Raven
	raven.Capture(packet, tags)
}

func runSentryWorker(ch chan map[string]interface{}) {
	log.Logger.Info("[sentry] push-worker is ready to receive message...")

	for {
		data := <-sentryChannel
		//metrics.GetSentryChannelUnreadGuage().Set(float64(len(sentryChannel)))
		//metrics.GetSentryChannelErrorCounter().Inc()
		sendSentryPacketV2(data)
	}
}

func startSentryWorker() {
	log.Logger.Info("[sentry] start push-worker")
	defer func() {
		if err := recover(); err != nil {
			stack := debug.Stack()
			fmt.Printf("\n>>>>>>>>>>>>>>>>>>>>\n%v\n%s\n<<<<<<<<<<<<<<<<<<<<\n", err, string(stack))
			//restart worker
			go startSentryWorker()
		}
	}()

	runSentryWorker(sentryChannel)
}

func init() {
	if isEnableSentry() {
		sentryDSN := config.ServiceConfig.String("sentry::SENTRY_DSN")
		raven.SetDSN(sentryDSN)
		log.Logger.Info(fmt.Sprintf("[sentry] use DSN:%s ", sentryDSN))
		go startSentryWorker()
	} else {
		log.Logger.Info("[sentry] sentry is DISABLED!!!")
	}
}
