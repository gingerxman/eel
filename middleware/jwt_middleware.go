package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/gingerxman/eel/config"
	"github.com/gingerxman/eel/handler"
	"github.com/gingerxman/eel/log"
	"strings"
)

var SALT string = "030e2cf548cf9da683e340371d1a74ee"

type JWTMiddleware struct {
	handler.Middleware
}

func (this *JWTMiddleware) ProcessRequest(ctx *handler.Context) {
	bCtx := ctx.GetBusinessContext()
	if ctx.Get("__shouldSkipAuthCheck").(bool) {
		if config.Runtime.NewBusinessContext != nil {
			bCtx = config.Runtime.NewBusinessContext(bCtx, ctx.Request.HttpRequest, 0, "", nil) //bCtx is for "business context"
		}
		ctx.SetBusinessContext(bCtx)
		return
	}
	
	//get jwt token
	jwtToken := ctx.Request.Header("AUTHORIZATION");
	if jwtToken == "" {
		//for dev
		jwtToken = ctx.Request.Query("_jwt")
	}
	
	if jwtToken != "" {
		items := strings.Split(jwtToken, ".")
		if len(items) != 3 {
			//jwt token 格式不对
			response := handler.MakeErrorResponse(500, "jwt:invalid_jwt_token", "无效的jwt token 1")
			ctx.Response.JSON(response)
			return
		}
		
		headerB64Code, payloadB64Code, expectedSignature := items[0], items[1], items[2]
		message := fmt.Sprintf("%s.%s", headerB64Code, payloadB64Code)
		
		h := hmac.New(sha256.New, []byte(SALT))
		h.Write([]byte(message))
		actualSignature := base64.StdEncoding.EncodeToString(h.Sum(nil));
		
		if expectedSignature != actualSignature {
			//jwt token的signature不匹配
			response := handler.MakeErrorResponse(500, "jwt:invalid_jwt_token", "无效的jwt token 2")
			ctx.Response.JSON(response)
			return
		}
		
		decodeBytes, err := base64.StdEncoding.DecodeString(payloadB64Code)
		if err != nil {
			log.Logger.Fatal(err)
		}
		js, err := simplejson.NewJson([]byte(decodeBytes))
		
		if err != nil {
			response := handler.MakeErrorResponse(500, "jwt:invalid_jwt_token", "无效的jwt token 3")
			ctx.Response.JSON(response)
			return
		}
		
		userId, err := js.Get("uid").Int()
		if err != nil {
			log.Logger.Fatal(err)
			response := handler.MakeErrorResponse(500, "jwt:invalid_jwt_token", "无效的jwt token 4")
			ctx.Response.JSON(response)
			return
		}
		
		if config.Runtime.NewBusinessContext != nil {
			bCtx = config.Runtime.NewBusinessContext(bCtx, ctx.Request.HttpRequest, userId, jwtToken, js) //bCtx is for "business context"
		}
		ctx.SetBusinessContext(bCtx)
	} else {
		response := handler.MakeErrorResponse(500, "jwt:invalid_jwt_token", "无效的jwt token 5")
		ctx.Response.JSON(response)
		return
	}
}

func (this *JWTMiddleware) ProcessResponse(ctx *handler.Context) {
	log.Logger.Info("i am in jwt middleware process response")
}

func init() {
}

