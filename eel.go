// Copyright 2018 eel Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eel

import (
	"context"
	"github.com/gingerxman/eel/paginate"
	"net/http"
	"time"

	"encoding/json"
	"fmt"

	"github.com/gingerxman/eel/config"
	"github.com/gingerxman/eel/db"
	"github.com/gingerxman/eel/handler"
	"github.com/gingerxman/eel/log"
	_ "github.com/gingerxman/eel/logo"
	"github.com/gingerxman/eel/middleware"
	"github.com/gingerxman/eel/rest_client"
	"github.com/gingerxman/eel/router"
	"github.com/gingerxman/eel/tracing"
	"github.com/gingerxman/eel/utils"
	"github.com/gingerxman/gorm"
	"go.uber.org/zap"
)

type Request handler.Request

var endRunning chan bool

//
//func handler(resp http.ResponseWriter, req *http.Request) {
//
//}

// export inner type
type Context = handler.Context
type RestResource = handler.RestResource
type Map = handler.Map
type BoolOption = map[string]bool
type Middleware = handler.Middleware
type Model = db.Model
type DeletableModel = db.DeletableModel
type FillOption = map[string]bool
type INextPageInfo = paginate.INextPageInfo
type PageInfo = paginate.PageInfo
type PaginateResult = paginate.PaginateResult

var Logger *zap.SugaredLogger = log.Logger
var Runtime = config.Runtime
var Tracer = tracing.Tracer
var Paginate = paginate.Paginate
var MockPaginate = paginate.MockPaginate

// export Middleware
type JWTMiddleware = middleware.JWTMiddleware

type IModel interface {
	TableName() string
}

type RepositoryBase struct {
	Ctx context.Context
}

type ServiceBase struct {
	Ctx context.Context
}

type EntityBase struct {
	Ctx   context.Context
	Model interface{}
}

func RegisterResource(resource handler.RestResourceInterface) {
	router.DoRegisterResource(resource)
}

func RegisterMiddleware(middleware handler.MiddlewareInterface) {
	router.DoRegisterMiddleware(middleware)
}

func RegisterModel(model interface{}) {
	db.RegisterModel(model)
}

func GetRegisteredModels() []interface{} {
	return db.GetRegisteredModels()
}

func NewBusinessError(code string, msg string) *utils.BusinessError {
	return utils.NewBusinessError(code, msg)
}

func MakeErrorResponse(code int32, errCode string, errMsg string, innerErrMsgs ...string) *handler.RestResponse {
	return MakeErrorResponse(code, errCode, errMsg, innerErrMsgs...)
}

func MakeResponse(data interface{}) *handler.RestResponse {
	return MakeResponse(data)
}

func NewResource(ctx context.Context) *rest_client.Resource {
	return rest_client.NewResource(ctx)
}

func ToJsonString(obj interface{}) string {
	bytes, _ := json.Marshal(obj)
	return string(bytes)
}

func GetOrmFromContext(ctx context.Context) *gorm.DB {
	o := ctx.Value("orm")
	return o.(*gorm.DB)
}

const SERVICE_MODE_REST = "rest"
const SERVICE_MODE_CRON = "cron"
const SERVICE_MODE_EVENT = "event"

func GetServiceMode() string {
	serviceMode := config.ServiceConfig.String("system::SERVICE_MODE")
	if serviceMode != SERVICE_MODE_CRON && serviceMode != SERVICE_MODE_REST && serviceMode != SERVICE_MODE_EVENT {
		panic(fmt.Sprintf("[CRITICAL] invalid service mode '%s'", serviceMode))
	}

	enableCronMode := config.ServiceConfig.DefaultBool("system::ENABLE_CRON_MODE", false)
	if enableCronMode {
		serviceMode = SERVICE_MODE_CRON
	}

	return serviceMode
}

type Service struct {
	Handler *router.RestResourceRegister
	Server  *http.Server
}

func NewService() *Service {
	app := &Service{
		Handler: router.NewRestResourceRegister(),
		Server:  &http.Server{},
	}
	return app
}

func (this *Service) run() {
	host := config.ServiceConfig.String("service::HOST")
	httpPort := config.ServiceConfig.String("service::HTTP_PORT")
	addr := fmt.Sprintf("%s:%s", host, httpPort)

	this.Server.Handler = this.Handler
	readTimeout := time.Duration(config.ServiceConfig.DefaultInt("service::READ_TIMEOUT", 30))
	writeTimeout := time.Duration(config.ServiceConfig.DefaultInt("service::WRITE_TIMEOUT", 30))
	readTimeout = 30
	writeTimeout = 30
	this.Server.ReadTimeout = readTimeout * time.Second
	this.Server.WriteTimeout = writeTimeout * time.Second
	this.Server.Addr = addr

	Logger.Infof("http server Running on http://%s\n", this.Server.Addr)
	if err := this.Server.ListenAndServe(); err != nil {
		Logger.Fatalf("ListenAndServe: ", err)
		time.Sleep(100 * time.Microsecond)
		endRunning <- true
	}
}

func RunService() {
	//fmt.Println(logo)
	service := NewService()
	endRunning = make(chan bool, 1)
	go func() {
		service.run()
	}()
	<-endRunning
}

func init() {
}
