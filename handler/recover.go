package handler

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gingerxman/eel/log"
	"github.com/gingerxman/eel/tracing"
	"github.com/gingerxman/eel/utils"
	"github.com/gingerxman/gorm"
	"github.com/opentracing/opentracing-go"
	"runtime"
	"runtime/debug"
	"strings"
)

func RecoverPanic(ctx *Context) {
	log.Logger.Debug("[router] in RecoverPanic...")
	cachedSpan := ctx.Get("rootSpan")
	var rootSpan opentracing.Span
	if cachedSpan != nil {
		rootSpan = cachedSpan.(opentracing.Span)
	}
	if err := recover(); err != nil {
		orm := ctx.Get("orm")
		if orm != nil {
			log.Logger.Info("[ORM] rollback transaction")
			var subSpan opentracing.Span
			if cachedSpan != nil {
				subSpan = tracing.CreateSubSpan(rootSpan, "db-rollback")
			}
			orm.(*gorm.DB).Rollback()
			if subSpan != nil {
				subSpan.Finish()
			}
		}

		errMsg := ""
		if be, ok := err.(*utils.BusinessError); ok {
			errMsg = fmt.Sprintf("%s:%s", be.ErrCode, be.ErrMsg)

		} else {
			errMsg = fmt.Sprintf("%s", err)
		}

		var buffer bytes.Buffer
		buffer.WriteString(fmt.Sprintf("[Unprocessed_Exception] %s\n", errMsg))
		buffer.WriteString(fmt.Sprintf("Request URL: %s\n", ctx.Request.URL()))
		for i := 1; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			buffer.WriteString(fmt.Sprintf("%s:%d\n", file, line))
		}
		log.Logger.Info(buffer.String())

		//记录到sentry
		{
			errMsg := ""
			if be, ok := err.(*utils.BusinessError); ok {
				errMsg = fmt.Sprintf("%s - %s", be.ErrCode, be.ErrMsg)
			} else {
				errMsg = fmt.Sprint(err)
			}
			CaptureErrorToSentry(ctx, errMsg)
		}

		if be, ok := err.(*utils.BusinessError); ok {
			ctx.Response.ErrorWithCode(500, be.ErrCode, be.ErrMsg, "")
		} else {
			ctx.Response.ErrorWithCode(531, "system:exception", fmt.Sprintf("%s", err), "")
		}
	} else {
		respCode := ctx.Response.ResponseWriter.Header().Get("X-Biz-Code")
		method := strings.ToLower(ctx.Request.HttpRequest.Method)
		if respCode == "500" && method == "post" {
			orm := ctx.Get("orm")
			if orm != nil {
				var subSpan opentracing.Span
				if cachedSpan != nil {
					subSpan = tracing.CreateSubSpan(rootSpan, "db-rollback")
				}
				orm.(*gorm.DB).Rollback()
				if subSpan != nil {
					subSpan.Finish()
				}
			}

			//记录到sentry
			{
				errMsg := fmt.Sprintf("%s - %s", ctx.Response.ErrCode, ctx.Response.ErrMsg)
				CaptureErrorToSentry(ctx, errMsg)
			}
		} else {
			orm := ctx.Get("orm")
			if orm != nil {
				log.Logger.Debug("[ORM] commit transaction")
				var subSpan opentracing.Span
				if cachedSpan != nil {
					subSpan = tracing.CreateSubSpan(rootSpan, "db-commit")
				}
				orm.(*gorm.DB).Commit()
				if subSpan != nil {
					subSpan.Finish()
				}

				ctx.EmitAll()
			}

			// 记录GET 500到sentry
			if respCode == "500" && method == "get" {
				errMsg := fmt.Sprintf("%s - %s", ctx.Response.ErrCode, ctx.Response.ErrMsg)
				CaptureErrorToSentry(ctx, errMsg)
			}
		}
	}
}

// RecoverFromCronTaskPanic crontask的recover
func RecoverFromCronTaskPanic(ctx context.Context) {
	o := ctx.Value("orm")
	if err := recover(); err != nil {
		log.Logger.Info("recover from cron task panic...")
		if o != nil {
			o.(*gorm.DB).Rollback()
			log.Logger.Warn("[ORM] rollback transaction for cron task")
		}

		log.Logger.Warn(string(debug.Stack()))
	}
}
