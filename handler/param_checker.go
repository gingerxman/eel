package handler

import (
	"github.com/bitly/go-simplejson"
	"strconv"
	"strings"
)

// CheckArgs: check whether param is right
func CheckArgs(r RestResourceInterface, ctx *Context) {
	req := ctx.Request
	method := req.Method()
	
	//if app, ok := r.AppController.(RestResourceInterface); ok {
	method2parameters := r.GetParameters()
	if method2parameters != nil {
		if parameters, ok := method2parameters[method]; ok {
			actualParams := req.Input()
			for _, param := range parameters {
				colonPos := strings.Index(param, ":")
				paramType := "string"
				if colonPos != -1 {
					paramType = param[colonPos+1 : len(param)]
					param = param[0:colonPos]
				}
				
				canMissParam := false
				if param[0] == '?' {
					canMissParam = true
					param = param[1:]
				}
				if _, ok := actualParams[param]; !ok {
					if !canMissParam {
						returnValidateParameterFailResponse(ctx, param, paramType, "no paramter provided")
						return
					} else {
						continue
					}
				}
				
				if paramType == "string" {
					//value := r.GetString(param)
				} else if paramType == "int" {
					_, err := req.GetInt64(param)
					if err != nil {
						returnValidateParameterFailResponse(ctx, param, paramType, err.Error())
						return
					} else {
						//requestData[param] = value
					}
				} else if paramType == "bool" {
					value := req.GetString(param)
					_, err := strconv.ParseBool(value)
					if err != nil {
						returnValidateParameterFailResponse(ctx, param, paramType, err.Error())
						return
					} else {
						//requestData[param] = result
					}
				} else if paramType == "json" {
					value := req.GetString(param)
					//						if value == "" && canMissParam == true {
					//							goto set_orm
					//						}
					js, err := simplejson.NewJson([]byte(value))
					if err != nil {
						returnValidateParameterFailResponse(ctx, param, paramType, err.Error())
						return
					} else {
						data, err := js.Map()
						if err != nil {
							returnValidateParameterFailResponse(ctx, param, paramType, err.Error())
							return
						} else {
							if param == "filters" {
								req.SetFilters(data)
							} else {
								req.SetJSON(param, data)
							}
						}
					}
				} else if paramType == "json-array" {
					value := req.GetString(param)
					js, err := simplejson.NewJson([]byte(value))
					if err != nil {
						returnValidateParameterFailResponse(ctx, param, paramType, err.Error())
						return
					} else {
						data, err := js.Array()
						if err != nil {
							returnValidateParameterFailResponse(ctx, param, paramType, err.Error())
							return
						} else {
							req.SetJSONArray(param, data)
						}
					}
				}
			}
			
			for key, value := range actualParams {
				if strings.HasPrefix(key, "__f") {
					req.SetFilter(key, value)
				}
			}
		}
	}
	//}
}
