package rest_client

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/gingerxman/eel/config"
	"github.com/gingerxman/eel/handler"
	"github.com/gingerxman/eel/log"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type ResourceResponse struct {
	RespData *simplejson.Json
}

func (this *ResourceResponse) IsSuccess() bool {
	code, _ := this.RespData.Get("code").Int()
	return code == 200
}

func (this *ResourceResponse) Data() *simplejson.Json {
	return this.RespData.Get("data")
}

func (this *ResourceResponse) ErrCode() string {
	if this.IsSuccess() {
		return ""
	} else {
		return this.RespData.Get("errCode").MustString()
	}
}

type Resource struct {
	Ctx            context.Context
	CustomJWTToken string
}

func (this *Resource) request(method string, service string, resource string, data handler.Map) (respData *ResourceResponse, err error) {
	var jwtToken string
	if this.CustomJWTToken != "" {
		jwtToken = this.CustomJWTToken
	} else {
		ijwt := this.Ctx.Value("jwt")
		if ijwt != nil {
			jwtToken = ijwt.(string)
		}
	}

	apiServerHost := config.ServiceConfig.String("api::API_SERVER_HOST")
	//创建client
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}

	//构建url.Values
	params := url.Values{"_v": {"1"}}

	//处理resource
	pos := strings.LastIndexByte(resource, '.')
	resource = fmt.Sprintf("%s/%s", resource[:pos], resource[pos+1:])

	//构建request
	apiUrl := fmt.Sprintf("http://%s/%s/%s/", apiServerHost, service, resource)
	fmt.Println(apiUrl)
	var req *http.Request
	if method == "GET" {
		for k, v := range data {
			value := ""
			switch t := v.(type) {
			case int:
				value = fmt.Sprintf("%d", v)
			case bool:
				value = fmt.Sprint("%t", v)
			case string:
				value = v.(string)
			default:
				log.Logger.Warn("unknown type: ", t)
			}
			params.Set(k, value)
		}
		apiUrl += "?" + params.Encode()
		log.Logger.Info("apiUrl: ", apiUrl)
		//strings.NewReader(values.Encode())

		req, err = http.NewRequest("GET", apiUrl, nil)
	} else {
		if method == "PUT" {
			params.Set("_method", "put")
		} else if method == "DELETE" {
			params.Set("_method", "delete")
		}
		apiUrl += "?" + params.Encode()
		log.Logger.Info("apiUrl: ", apiUrl)

		values := url.Values{}
		for k, v := range data {
			value := ""
			switch t := v.(type) {
			case int:
				value = fmt.Sprintf("%d", v)
			case bool:
				value = fmt.Sprint("%t", v)
			case string:
				value = v.(string)
			default:
				log.Logger.Warn("unknown type: ", t)
			}
			values.Set(k, value)
		}

		req, err = http.NewRequest("POST", apiUrl, strings.NewReader(values.Encode()))
	}
	if err != nil {
		return nil, err
	}

	if method != "GET" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if jwtToken != "" {
		req.Header.Set("AUTHORIZATION", jwtToken)
	}

	//inject open tracing
	span := opentracing.SpanFromContext(this.Ctx)
	if span != nil {
		ext.SpanKindRPCClient.Set(span)
		ext.HTTPUrl.Set(span, apiUrl)
		ext.HTTPMethod.Set(span, method)
		span.Tracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header),
		)
	}

	//执行request，获得response
	resp, err := netClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	//获取response的内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	jsonObj := new(simplejson.Json)
	err = jsonObj.UnmarshalJSON(body)
	if err != nil {
		return nil, err
	}

	resourceResp := new(ResourceResponse)
	resourceResp.RespData = jsonObj
	//fmt.Println(string(body))

	if resourceResp.IsSuccess() {
		return resourceResp, nil
	} else {
		log.Logger.Error(jsonObj)
		return resourceResp, errors.New("business_error")
	}
}

func (this *Resource) Get(service string, resource string, data handler.Map) (resp *ResourceResponse, err error) {
	return this.request("GET", service, resource, data)
}

func (this *Resource) Put(service string, resource string, data handler.Map) (resp *ResourceResponse, err error) {
	return this.request("PUT", service, resource, data)
}

func (this *Resource) Post(service string, resource string, data handler.Map) (resp *ResourceResponse, err error) {
	return this.request("POST", service, resource, data)
}

func (this *Resource) Delete(service string, resource string, data handler.Map) (resp *ResourceResponse, err error) {
	return this.request("DELETE", service, resource, data)
}

func (this *Resource) LoginAs(username string) *Resource {
	password := config.ServiceConfig.String("system::SUPER_PASSWORD")
	accountService := config.ServiceConfig.DefaultString("system::ACCOUNT_SERVICE", "ginger-account")
	resp, err := this.Put(accountService, "login.logined_corp_user", handler.Map{
		"username": username,
		"password": password,
	})
	if err != nil {
		log.Logger.Error(err)
		return nil
	}

	respData := resp.Data()
	this.CustomJWTToken, _ = respData.Get("jwt").String()
	return this
}

func (this *Resource) LoginAsManager() *Resource {
	managerAccount := config.ServiceConfig.DefaultString("system::MANAGER_ACCOUNT", "ginger")
	return this.LoginAs(managerAccount)
}

func NewResource(ctx context.Context) *Resource {
	resource := new(Resource)
	resource.Ctx = ctx
	return resource
}
