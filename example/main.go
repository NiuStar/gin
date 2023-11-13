package main

import (
	"github.com/NiuStar/gin"
)

type TestRequest struct {
	UserName string `json:"user_name" 备注:"用户名" 默认值:"admin"`
	Password string `json:"password" 备注:"密码" 默认值:"123456"`
}
type TestResponse struct {
	UserId       string       `json:"user_id" 备注:"用户id" 默认值:"1"`
	UserName     string       `json:"user_name" 备注:"用户名" 默认值:"admin"`
	CreateTime   string       `json:"create_time" 备注:"创建时间" 默认值:"2020-01-01 00:00:00"`
	Test1Request Test1Request `json:"test1_request" 备注:"用户名" 默认值:"admin"`
}

type Test1Request struct {
	UserName1    string          `json:"user_name" 备注:"用户名" 默认值:"admin"`
	Test2Request []*Test2Request `json:"test2_request" 备注:"用户名" 默认值:"admin"`
}

type Test2Request struct {
	UserName2 string `json:"user_name" 备注:"用户名" 默认值:"admin"`
}

func Login(c *gin.Context) {
	var request TestRequest
	c.BindJSON(&request)
}
func main() {
	r := gin.Default("测试登录项目3")
	apiGroup := r.Group("/api", nil, nil)
	testGroup := apiGroup.Group("/test", nil, nil)
	testGroup.Handle("POST", "/login", &TestRequest{}, &TestResponse{}, Login)
	r.RunWithDoc("test", "123456", "https://showdoc.ai00.xyz", ":8082")
}
