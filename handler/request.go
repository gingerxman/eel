package handler

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"github.com/gingerxman/eel/paginate"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

//Request
type Request struct {
	HttpRequest *http.Request
	Name2JSON map[string]map[string]interface{}
	Name2JSONArray map[string][]interface{}
	Filters map[string]interface{}
}

func (r *Request) Reset(request *http.Request) {
	r.HttpRequest = request
	r.Name2JSONArray = nil
	r.Name2JSON = nil
	r.Filters = nil
}

// Query returns input data item string by a given string.
func (r *Request) Query(key string) string {
	//if val := r.Param(key); val != "" {
	//	return val
	//}
	if r.HttpRequest.Form == nil {
		r.HttpRequest.ParseForm()
	}
	return r.HttpRequest.Form.Get(key)
}

func (r *Request) Input() url.Values {
	if r.HttpRequest.Form == nil {
		r.HttpRequest.ParseForm()
	}
	
	return r.HttpRequest.Form
}

func (r *Request) SetJSON(key string, value map[string]interface{}) {
	if r.Name2JSON == nil {
		r.Name2JSON = make(map[string]map[string]interface{})
	}
	r.Name2JSON[key] = value
}

func (r *Request) SetJSONArray(key string, value []interface{}) {
	if r.Name2JSONArray == nil {
		r.Name2JSONArray = make(map[string][]interface{})
	}
	r.Name2JSONArray[key] = value
}

func (r *Request) SetFilter(key string, value interface{}) {
	if r.Filters == nil {
		r.Filters = make(map[string]interface{})
	}
	r.Filters[key] = value
}

func (r *Request) SetFilters(filters map[string]interface{}) {
	r.Filters = filters
}

//GetJSONArray 与key对应的返回json array数据
func (r *Request) GetJSONArray(key string) []interface{} {
	if data, ok := r.Name2JSONArray[key]; ok {
		return data
	} else {
		return nil
	}
}

func (r *Request) GetIntArray(key string) []int {
	values := make([]int, 0)
	if datas, ok := r.Name2JSONArray[key]; ok {
		for _, data := range datas {
			intValue, _ := strconv.Atoi(data.(json.Number).String())
			values = append(values, intValue)
		}
		return values
	} else {
		return nil
	}
}

func (r *Request) GetStringArray(key string) []string {
	values := make([]string, 0)
	if datas, ok := r.Name2JSONArray[key]; ok {
		for _, data := range datas {
			values = append(values, data.(string))
		}
		return values
	} else {
		return nil
	}
}

func (r *Request) GetBoolOptions(key string) map[string]bool {
	boolOption := make(map[string]bool)
	if data, ok := r.Name2JSON[key]; ok {
		for key, _ := range data {
			boolOption[key] = true
		}
	} else {
		return nil
	}
	
	return boolOption
}

//GetJSONArray 与key对应的返回json map数据
func (r *Request) GetJSON(key string) map[string]interface{} {
	if data, ok := r.Name2JSON[key]; ok {
		return data
	} else {
		return nil
	}
}

func (r *Request) GetFillOptions(key string) FillOption {
	fillOption := FillOption{}
	if data, ok := r.Name2JSON[key]; ok {
		for key, _ := range data {
			fillOption[key] = true
		}
	} else {
		return fillOption
	}
	
	return fillOption
}

func (r *Request) GetFilters() map[string]interface{} {
	return r.Filters
}

func convertToOrmFilter(filters map[string]interface{}) map[string]interface{} {
	params := make(map[string]interface{})
	var key string
	var match string
	
	for k, v := range filters {
		if !(strings.HasPrefix(k, "__f")) {
			params[k] = v
		} else {
			keyString := strings.Split(k, "-")
			if len(keyString) == 3 {
				switch keyString[2] {
				case "equal":
					match = "exact"
				case "contain":
					match = "contains"
				case "gt", "gte", "lt", "lte", "in":
					match = keyString[2]
				case "range":
					val := v.([]interface{})
					start, _ := val[0].(json.Number).Int64()
					stop, _ := val[1].(json.Number).Int64()
					
					values := make([]int64, 0)
					for i := start; i < stop; i++ {
						values = append(values, i)
					}
					
					match = "in"
					v = values
				}
				key = fmt.Sprintf("%s__%s", keyString[1], match)
			} else {
				key = fmt.Sprintf("%s", keyString[1])
			}
			params[key] = v
		}
	}
	return params
}

func (r *Request) GetOrmFilters() map[string]interface{} {
	return convertToOrmFilter(r.GetFilters())
}

func (r *Request) GetString(key string, def ...string) string {
	if v := r.Query(key); v != "" {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

func (r *Request) Method() string {
	method := r.HttpRequest.Method
	if method == "POST" {
		input := r.Input()
		_method := input.Get("_method")
		if _method == "put" {
			method = "PUT"
		} else if _method == "delete" {
			method = "DELETE"
		}
	}
	return method
}

func (r *Request) RawMethod() string {
	return r.HttpRequest.Method
}

// GetInt returns input as an int or the default value while it's present and input is blank
func (r *Request) GetInt(key string, def ...int) (int, error) {
	strv := r.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.Atoi(strv)
}

// GetInt8 return input as an int8 or the default value while it's present and input is blank
func (r *Request) GetInt8(key string, def ...int8) (int8, error) {
	strv := r.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	i64, err := strconv.ParseInt(strv, 10, 8)
	return int8(i64), err
}

// GetUint8 return input as an uint8 or the default value while it's present and input is blank
func (r *Request) GetUint8(key string, def ...uint8) (uint8, error) {
	strv := r.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	u64, err := strconv.ParseUint(strv, 10, 8)
	return uint8(u64), err
}

// GetInt16 returns input as an int16 or the default value while it's present and input is blank
func (r *Request) GetInt16(key string, def ...int16) (int16, error) {
	strv := r.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	i64, err := strconv.ParseInt(strv, 10, 16)
	return int16(i64), err
}

// GetUint16 returns input as an uint16 or the default value while it's present and input is blank
func (r *Request) GetUint16(key string, def ...uint16) (uint16, error) {
	strv := r.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	u64, err := strconv.ParseUint(strv, 10, 16)
	return uint16(u64), err
}

// GetInt32 returns input as an int32 or the default value while it's present and input is blank
func (r *Request) GetInt32(key string, def ...int32) (int32, error) {
	strv := r.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	i64, err := strconv.ParseInt(strv, 10, 32)
	return int32(i64), err
}

// GetUint32 returns input as an uint32 or the default value while it's present and input is blank
func (r *Request) GetUint32(key string, def ...uint32) (uint32, error) {
	strv := r.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	u64, err := strconv.ParseUint(strv, 10, 32)
	return uint32(u64), err
}

// GetInt64 returns input value as int64 or the default value while it's present and input is blank.
func (r *Request) GetInt64(key string, def ...int64) (int64, error) {
	strv := r.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.ParseInt(strv, 10, 64)
}

// GetUint64 returns input value as uint64 or the default value while it's present and input is blank.
func (r *Request) GetUint64(key string, def ...uint64) (uint64, error) {
	strv := r.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.ParseUint(strv, 10, 64)
}

// GetBool returns input value as bool or the default value while it's present and input is blank.
func (r *Request) GetBool(key string, def ...bool) (bool, error) {
	strv := r.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.ParseBool(strv)
}

// GetFloat returns input value as float64 or the default value while it's present and input is blank.
func (r *Request) GetFloat(key string, def ...float64) (float64, error) {
	strv := r.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.ParseFloat(strv, 64)
}

// URL returns request url path (without query string, fragment).
func (r *Request) URL() string {
	return r.HttpRequest.URL.Path
}

func (r *Request) Header(key string) string {
	return r.HttpRequest.Header.Get(key)
}

func (r *Request) getInt(key string, def ...int) (int, error) {
	strv := r.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.Atoi(strv)
}

func (r *Request) GetFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return r.HttpRequest.FormFile(key)
}

func (r *Request) GetPageInfo() *paginate.PageInfo {
	fromParam := r.Query("_p_from")
	if fromParam != "" {
		fromId, _ := r.getInt("_p_from")
		countPerPage, _ := r.getInt("_p_count", 20)
		return &paginate.PageInfo{
			Page:         -1,
			FromId:       fromId,
			CountPerPage: countPerPage,
			Mode:         "apiserver",
			Direction:    "asc",
		}
	} else {
		page, _ := r.getInt("page", 1)
		countPerPage, _ := r.getInt("count_per_page", 20)
		return &paginate.PageInfo{
			Page:         page,
			FromId:       0,
			CountPerPage: countPerPage,
			Mode:         "backend",
			Direction:    "asc",
		}
	}
}