package showdoc

import (
	http2 "alicode.mukj.cn/yjkj.ink/work/http"
	"errors"
	"fmt"
	"net/http"
	"sync"
)

type ShowDoc struct {
	UserName string
	Password string
	Host     string
	Header   http.Header
	ApiKey   *ApiKey
}

type Response struct {
	ErrorCode int64 `json:"error_code"`
}

var instance *ShowDoc
var once sync.Once

func Instance() *ShowDoc {
	once.Do(func() {
		instance = &ShowDoc{}
	})
	return instance
}

func (doc *ShowDoc) Login(userName, password, host string) error {
	user := map[string]string{"username": userName, "password": password}
	doc.Host = host
	resp := http2.POSTFormData(fmt.Sprintf("%s%s", host, Login), user)
	if resp.Error() != nil {
		fmt.Println("登录showdoc失败1,", resp.Error())
		return resp.Error()
	}
	var result *Response
	err := resp.Resp(&result)
	if err != nil {
		fmt.Println("登录showdoc失败2,", string(resp.Byte()))
		return errors.New(string(resp.Byte()))
	}
	if result.ErrorCode == 0 {
		doc.Header = resp.Header()
		doc.UserName = userName
		doc.Password = password
		return nil
	}
	fmt.Println("登录showdoc失败3,", string(resp.Byte()))
	return errors.New(string(resp.Byte()))
}

func (doc *ShowDoc) GetServiceList() *ItemResponse {
	if doc.Header == nil {
		return nil
	}
	resp := http2.GETWithHeader(fmt.Sprintf("%s%s", doc.Host, GetServiceList), nil, &doc.Header)
	if resp.Error() != nil {
		fmt.Println("获取文档列表失败1,", resp.Error())
		return nil
	}
	var result *ItemResponse
	err := resp.Resp(&result)
	if err != nil {
		fmt.Println("获取文档列表失败2,", string(resp.Byte()))
		return nil
	}
	if result.ErrorCode == 0 {

		return result
	}
	fmt.Println("获取文档列表失败3,", string(resp.Byte()))
	return nil
}

func (doc *ShowDoc) GetApiKey(itemId string) *ApiKeyResponse {
	if doc.Header == nil {
		return nil
	}
	item := map[string]string{"item_id": itemId}
	resp := http2.POSTFormDataWithHeader(fmt.Sprintf("%s%s", doc.Host, GetApiKey), item, &doc.Header)
	if resp.Error() != nil {
		fmt.Println("获取ApiKey失败1,", resp.Error())
		return nil
	}
	var result *ApiKeyResponse
	err := resp.Resp(&result)
	if err != nil {
		fmt.Println("获取ApiKey失败2,", string(resp.Byte()))
		return nil
	}
	if result.ErrorCode == 0 {

		return result
	}
	fmt.Println("获取ApiKey失败3,", string(resp.Byte()))
	return nil
}

func (doc *ShowDoc) AddService(name string) error {
	if doc.Header == nil {
		return nil
	}
	item := map[string]string{"item_type": "1", "item_name": name, "item_description": "自动创建的", "item_domain": ""}
	resp := http2.POSTFormDataWithHeader(fmt.Sprintf("%s%s", doc.Host, AddService), item, &doc.Header)
	if resp.Error() != nil {
		fmt.Println("添加项目文档失败1,", resp.Error())
		return resp.Error()
	}
	var result *Response
	err := resp.Resp(&result)
	if err != nil {
		fmt.Println("添加项目文档失败2,", string(resp.Byte()))
		return errors.New(string(resp.Byte()))
	}
	if result.ErrorCode == 0 {

		return nil
	}
	fmt.Println("添加项目文档失败3,", string(resp.Byte()))
	return errors.New(string(resp.Byte()))
}

func (doc *ShowDoc) DeleteService(itemId string) error {
	if doc.Header == nil {
		return nil
	}
	item := map[string]string{"item_id": itemId, "password": doc.Password}
	resp := http2.POSTFormDataWithHeader(fmt.Sprintf("%s%s", doc.Host, DelService), item, &doc.Header)
	if resp.Error() != nil {
		fmt.Println("删除项目文档失败1,", resp.Error())
		return resp.Error()
	}
	var result *Response
	err := resp.Resp(&result)
	if err != nil {
		fmt.Println("删除项目文档失败2,", string(resp.Byte()))
		return errors.New(string(resp.Byte()))
	}
	if result.ErrorCode == 0 {

		return nil
	}
	fmt.Println("删除项目文档失败3,", string(resp.Byte()))
	return errors.New(string(resp.Byte()))
}
