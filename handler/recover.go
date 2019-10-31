package handler

import (
	"fmt"
	"bytes"
	"runtime"
	"github.com/gingerxman/eel/utils"
	"github.com/gingerxman/eel/log"
	"github.com/gingerxman/gorm"
	"github.com/gingerxman/eel/tracing"
	"github.com/opentracing/opentracing-go"
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
			log.Logger.Warn("[ORM] rollback transaction")
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
		log.Logger.Error(buffer.String())
		
		if be, ok := err.(*utils.BusinessError); ok {
			ctx.Response.ErrorWithCode(500, be.ErrCode, be.ErrMsg, "")
		} else {
			ctx.Response.ErrorWithCode(531, "system:exception", fmt.Sprintf("%s", err), "")
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
		}
	}
}
