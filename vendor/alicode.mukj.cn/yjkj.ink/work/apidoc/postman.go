package apidoc

import (
	"alicode.mukj.cn/yjkj.ink/work/utils/uuid"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type Info struct {
	PostmanID   string `json:"_postman_id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Schema      string `json:"schema"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type Url struct {
	Query    []*KeyValue `json:"query,omitempty"`
	Host     []string    `json:"host"`
	Path     []string    `json:"path"`
	Port     string      `json:"port,omitempty"`
	Protocol string      `json:"protocol"`
	Raw      string      `json:"raw"`
}

type Options struct {
	Raw struct {
		Language string `json:"language"`
	} `json:"raw"`
}

type Body struct {
	Mode       string   `json:"mode"`
	Options    *Options `json:"options,omitempty"`
	Raw        string   `json:"raw,omitempty"`
	UrlEncoded []*Value `json:"urlencoded,omitempty"`
	FormData   []*Value `json:"formdata,omitempty"`
}

type Value struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Request struct {
	Body   *Body    `json:"body,omitempty"`
	Header []*Value `json:"header"`
	Method string   `json:"method"`
	URL    Url      `json:"url"`
}

type Item struct {
	Item     []*Item       `json:"item,omitempty"`
	Name     string        `json:"name"`
	Request  *Request      `json:"request,omitempty"`
	Response []interface{} `json:"response,omitempty"`
}

func (item2 *Item) getFloderByName(name string) *Item {
	for _, item := range item2.Item {
		if item.Name == name && item.Request == nil {
			return item
		}
	}
	return nil
}

type PostmanDoc struct {
	Info                    Info     `json:"info"`
	Item                    []*Item  `json:"item"`
	ProtocolProfileBehavior struct{} `json:"protocolProfileBehavior"`
}

const MaxLevel = 4

func NewPostmanDoc(projectName string) *PostmanDoc {
	uuid, _ := uuid.NewV4()
	return &PostmanDoc{Info: Info{Name: projectName, Schema: "https://schema.getpostman.com/json/collection/v2.1.0/collection.json", PostmanID: uuid.String()}}
}

type ContentType string

const (
	NONE       ContentType = "none"
	FormData   ContentType = "formdata"
	UrlEncoded ContentType = "urlencoded"
	Json       ContentType = "raw"
)

func (doc *PostmanDoc) getFloderByName(name string) *Item {
	for _, item := range doc.Item {
		if item.Name == name && item.Request == nil {
			return item
		}
	}
	return nil
}

func (doc *PostmanDoc) WriteItem(apiName, flodrName, method, uri string, contentType ContentType, urlParams, request interface{}) {

	item := &Item{}
	item.Name = apiName
	item.Request = &Request{}
	item.Request.Method = method
	item.Request.URL.Raw = uri
	ruri, err := url.Parse(uri)
	if err == nil {
		item.Request.URL.Protocol = ruri.Scheme
		item.Request.URL.Host = strings.Split(strings.Split(ruri.Host, ":")[0], ".")
		item.Request.URL.Port = ruri.Port()
		item.Request.URL.Path = strings.Split(ruri.Path, "/")

	}

	if strings.ToUpper(method) == "POST" {
		if contentType == FormData {
			item.Request.Header = append(item.Request.Header, &Value{Key: "Content-Type", Value: "application/formdata", Type: "text"})
			item.Request.Body = createFormDataBody(request)
		} else if contentType == UrlEncoded {
			item.Request.Header = append(item.Request.Header, &Value{Key: "Content-Type", Value: "application/x-www-form-urlencoded", Type: "text"})
			item.Request.Body = createUrlEncodeBody(request)
		} else if contentType == Json {
			item.Request.Header = append(item.Request.Header, &Value{Key: "Content-Type", Value: "application/json", Type: "text"})
			item.Request.Body = createRawBody(request)
		}
	} else if request != nil {
		reqType := reflect.TypeOf(request)
		for reqType.Kind() == reflect.Ptr {
			reqType = reqType.Elem()
		}
		for i := 0; i < reqType.NumField(); i++ {
			value := getValueFromField(reqType.Field(i))
			item.Request.URL.Query = append(item.Request.URL.Query, &KeyValue{Key: value.Key, Value: value.Value})
		}
	}
	if values, err := url.ParseQuery(ruri.RawQuery); err == nil {
		for key, v := range values {
			for _, value := range v {
				item.Request.URL.Query = append(item.Request.URL.Query, &KeyValue{Key: key, Value: value})
			}
		}
	}
	if urlParams != nil {
		reqType := reflect.TypeOf(urlParams)
		for reqType.Kind() == reflect.Ptr {
			reqType = reqType.Elem()
		}
		for i := 0; i < reqType.NumField(); i++ {
			value := getValueFromField(reqType.Field(i))
			item.Request.URL.Query = append(item.Request.URL.Query, &KeyValue{Key: value.Key, Value: value.Value})
		}
	}
	if len(flodrName) > 0 {
		var item2 *Item
		list := strings.Split(flodrName, "/")
		pos := 0
		for _, v := range list {
			if v != "" {
				if pos == 0 {
					item2 = doc.getFloderByName(v)
					if item2 == nil {
						item2 = &Item{Name: v}
						doc.Item = append(doc.Item, item2)
					}
				} else {
					item3 := item2.getFloderByName(v)
					if item3 == nil {
						item3 = &Item{Name: v}
						item2.Item = append(item2.Item, item3)
					}
					item2 = item3
				}
				pos++
			}
		}
		if item2 == nil {
			doc.Item = append(doc.Item, item)
		} else {
			item2.Item = append(item2.Item, item)
		}
	} else {
		doc.Item = append(doc.Item, item)
	}

}

func getValueFromField(field reflect.StructField) *Value {
	value := &Value{}
	value.Type = "text"
	jsonName := field.Tag.Get("json")
	if len(jsonName) <= 0 {
		jsonName = field.Name
	}
	value.Key = jsonName
	value.Value = field.Tag.Get("default")
	if field.Type.Kind() == reflect.Int || field.Type.Kind() == reflect.Int32 || field.Type.Kind() == reflect.Int64 ||
		field.Type.Kind() == reflect.Int8 || field.Type.Kind() == reflect.Uint || field.Type.Kind() == reflect.Uint32 ||
		field.Type.Kind() == reflect.Uint64 || field.Type.Kind() == reflect.Uint8 || field.Type.Kind() == reflect.Float64 ||
		field.Type.Kind() == reflect.Float32 {

		if len(value.Value) <= 0 {
			value.Value = "0"
		}
	} else if field.Type.Kind() == reflect.Bool {
		if len(value.Value) <= 0 {
			value.Value = "false"
		}
	} else if field.Type.Kind() == reflect.Ptr {
		if field.Type.String() == "*os.File" {
			value.Type = "file"
		} else {
			v := getRawJsonFromType(field.Type, 0).Interface()
			b, _ := json.Marshal(v)
			value.Value = string(b)
		}
	} else if field.Type.Kind() == reflect.Struct {
		if field.Type.String() == "os.File" {
			value.Type = "file"
		} else {
			v := getRawJsonFromType(field.Type, 0).Interface()
			b, _ := json.Marshal(v)
			value.Value = string(b)
		}
	}
	return value
}

func getRawJsonFromType(reqType reflect.Type, num int) reflect.Value {
	if num >= MaxLevel {
		return reflect.ValueOf(nil)
	}
	num++
	//fmt.Println("reqType", reqType)

	isPtr := false
	oReqType := reqType
	if oReqType.Kind() == reflect.Ptr {
		oReqType = oReqType.Elem()
		isPtr = true
	}
	if oReqType.Kind() == reflect.Struct {
		req := reflect.New(oReqType).Elem()
		//fmt.Println("req", req.Type())
		/*if isPtr {
			req = req.Elem()
			fmt.Println("req2",req.Type())

		}*/
		for i := 0; i < oReqType.NumField(); i++ {
			name := oReqType.Field(i).Name
			if []byte(name)[0] >= 'a' && []byte(name)[0] <= 'z' {
				continue
			}
			if oReqType.Field(i).Type.Kind() == reflect.Int || oReqType.Field(i).Type.Kind() == reflect.Int32 || oReqType.Field(i).Type.Kind() == reflect.Int64 ||
				oReqType.Field(i).Type.Kind() == reflect.Int8 {

				v := oReqType.Field(i).Tag.Get("default")
				req.Field(i).SetInt(0)
				if len(v) > 0 {
					if v1, err := strconv.ParseInt(v, 10, 64); err == nil {
						req.Field(i).SetInt(v1)
					}
				}

			} else if oReqType.Field(i).Type.Kind() == reflect.Uint || oReqType.Field(i).Type.Kind() == reflect.Uint32 ||
				oReqType.Field(i).Type.Kind() == reflect.Uint64 || oReqType.Field(i).Type.Kind() == reflect.Uint8 {

				kind := oReqType.Field(i).Type.Kind()
				fmt.Println("kind", kind, oReqType)
				v := oReqType.Field(i).Tag.Get("default")
				req.Field(i).SetUint(uint64(0))
				if len(v) > 0 {
					if v1, err := strconv.ParseInt(v, 10, 64); err == nil {
						req.Field(i).SetUint(uint64(v1))
					}
				}

			} else if oReqType.Field(i).Type.Kind() == reflect.Float64 ||
				oReqType.Field(i).Type.Kind() == reflect.Float32 {
				v := oReqType.Field(i).Tag.Get("default")
				req.Field(i).SetFloat(0)
				if len(v) > 0 {
					if v1, err := strconv.ParseFloat(v, 64); err == nil {
						req.Field(i).SetFloat(v1)
					}
				}
			} else if oReqType.Field(i).Type.Kind() == reflect.Bool {
				v := oReqType.Field(i).Tag.Get("default")
				if v == "true" {
					req.Field(i).SetBool(true)
				} else {
					req.Field(i).SetBool(false)
				}
			} else if oReqType.Field(i).Type.Kind() == reflect.String {
				v := oReqType.Field(i).Tag.Get("default")
				req.Field(i).SetString(v)
			} else if oReqType.Field(i).Type.Kind() == reflect.Slice {
				//fmt.Println("reqType.Field(i).Type", oReqType.Field(i).Type)

				//fmt.Println("reqType.Field(i).Type.Elem()", oReqType.Field(i).Type.Elem())
				rs := reflect.MakeSlice(oReqType.Field(i).Type, 0, 0)
				if num < MaxLevel {
					v := getRawJsonFromType(oReqType.Field(i).Type.Elem(), num)
					//fmt.Println("v.Type", v.Type())
					if v.Kind() == reflect.Ptr {
						if v.IsValid() && !v.IsNil() && !rs.IsZero() && !v.IsZero() {
							rs = reflect.Append(rs, v)
						}
					} else {

						if v.IsValid() {
							name1 := v.Kind().String()
							fmt.Println(name1)
							rs = reflect.Append(rs, v)
							if oReqType.Field(i).Type.Elem().Kind() == reflect.Ptr && v.CanAddr() {
								//fmt.Println("v.Addr().Type", v.Addr().Type())
								//fmt.Println("rs.Type", rs.Type())
								//rs = reflect.AppendSlice(rs, v.Addr())

							} else if oReqType.Field(i).Type.Elem().Kind() != reflect.Ptr {

							}
						}

					}
				}

				req.Field(i).Set(rs)
			} else if oReqType.Field(i).Type.Kind() == reflect.Ptr {
				if num < MaxLevel {
					v := getRawJsonFromType(oReqType.Field(i).Type, num)
					if v.IsValid() && !v.IsNil() {
						req.Field(i).Set(v)
					}

					if v.CanAddr() {
						//	req.Field(i).Set(v.Addr())
					} else {

					}
				}
			} else if oReqType.Field(i).Type.Kind() == reflect.Struct {
				if num < MaxLevel {
					req.Field(i).Set(getRawJsonFromType(oReqType.Field(i).Type, num))
				}
			}
		}
		if isPtr {
			return req.Addr()
		} else {
			return req
		}
	} else if oReqType.Kind() == reflect.Ptr {
		if num < MaxLevel {
			return getRawJsonFromType(oReqType.Elem(), num).Addr()
		} else {
			return reflect.ValueOf(nil)
		}
	} else if reqType.Kind() == reflect.Int || reqType.Kind() == reflect.Int32 || reqType.Kind() == reflect.Int64 ||
		reqType.Kind() == reflect.Int8 {
		v := reflect.New(reqType).Elem()
		v.SetInt(0)
		if isPtr {
			return v.Addr()
		} else {
			return v
		}
	} else if reqType.Kind() == reflect.Uint || reqType.Kind() == reflect.Uint32 ||
		reqType.Kind() == reflect.Uint64 || reqType.Kind() == reflect.Uint8 {

		v := reflect.New(reqType).Elem()
		v.SetUint(0)
		if isPtr {
			return v.Addr()
		} else {
			return v
		}

	} else if reqType.Kind() == reflect.Float64 ||
		reqType.Kind() == reflect.Float32 {
		v := reflect.New(reqType).Elem()
		v.SetFloat(0)
		if isPtr {
			return v.Addr()
		} else {
			return v
		}
	} else if reqType.Kind() == reflect.Bool {
		v := reflect.New(reqType).Elem()
		v.SetBool(false)
		if isPtr {
			return v.Addr()
		} else {
			return v
		}
	} else if reqType.Kind() == reflect.String {
		v := reflect.New(reqType).Elem()
		v.SetString("")
		if isPtr {
			return v.Addr()
		} else {
			return v
		}
	}
	return reflect.ValueOf(nil)
}

func createFormDataBody(request interface{}) *Body {
	if request == nil {
		return nil
	}
	body := &Body{Mode: string(FormData)}

	reqType := reflect.TypeOf(request)
	for reqType.Kind() == reflect.Ptr {
		reqType = reqType.Elem()
	}
	for i := 0; i < reqType.NumField(); i++ {
		body.FormData = append(body.FormData, getValueFromField(reqType.Field(i)))
	}
	return body
}

func createUrlEncodeBody(request interface{}) *Body {
	if request == nil {
		return nil
	}
	body := &Body{Mode: string(UrlEncoded)}
	reqType := reflect.TypeOf(request)
	for reqType.Kind() == reflect.Ptr {
		reqType = reqType.Elem()
	}
	for i := 0; i < reqType.NumField(); i++ {
		body.UrlEncoded = append(body.UrlEncoded, getValueFromField(reqType.Field(i)))
	}
	return body
}

func createRawBody(request interface{}) *Body {

	body := &Body{Mode: string(Json)}
	body.Options = &Options{}
	body.Options.Raw.Language = "json"
	if request != nil {
		v := getRawJsonFromType(reflect.TypeOf(request), 0).Interface()
		b, _ := json.MarshalIndent(v, "", "\t")
		body.Raw = string(b)
	} else {
		body.Raw = "{}"
	}

	return body
}

func (doc *PostmanDoc) Save(fileName string) {
	data, _ := json.Marshal(doc)
	err := ioutil.WriteFile(fileName, data, 0777)
	if err != nil {
		fmt.Println("postman write err", err)
	}
}
