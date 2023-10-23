package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"time"
	//jsoniter "github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
)

type HttpResponse struct {
	result     []byte
	err        error
	header     http.Header
	StatusCode int
}

func (resp *HttpResponse) Resp(result interface{}) error {

	l := len(resp.result)
	if l > 0 && resp.result[0] == '<' {
		return xml.Unmarshal(resp.result, result)
	} else if l > 0 && (resp.result[0] == '{' || resp.result[0] == '[') {
		return json.Unmarshal(resp.result, result)
	}

	return errors.New("HttpResponse Resp数据字段返回异常")
}

func (resp *HttpResponse) Error() error {
	return resp.err
}

func (resp *HttpResponse) Header() http.Header {
	return resp.header
}

func (resp *HttpResponse) Byte() []byte {
	return resp.result
}

func POSTData(uri string, params url.Values) (resp *HttpResponse) {

	return Request("POST", uri, "", params.Encode())

}

func PATCHData(uri string, params url.Values) (resp *HttpResponse) {
	return Request("PATCH", uri, "", params.Encode())
}

func PATCHJson(uri string, params interface{}) (resp *HttpResponse) {
	return RequestJsonBody("PATCH", uri, params)
}

func POSTJson(uri string, params interface{}) (resp *HttpResponse) {
	return RequestJsonBody("POST", uri, params)
}

func RequestJsonBody(method, uri string, params interface{}) (resp *HttpResponse) {
	payload, err := makePayload(params)
	if err != nil {
		return &HttpResponse{err: err}
	}
	return Request(method, uri, "application/json", payload)
}

func PATCHFormDatan(uri string, postData map[string]string) (resp *HttpResponse) {
	return RequestFormData("PATCH", uri, postData)
}

func POSTFormData(uri string, postData map[string]string) (resp *HttpResponse) {
	return RequestFormData("POST", uri, postData)
}

func POSTFormUrlDecode(uri string, postData map[string]string) (resp *HttpResponse) {
	value := url.Values{}
	for key, d := range postData {
		value.Add(key, d)
	}
	return Request("POST", uri, "application/x-www-form-urlencoded", value.Encode())
}

func RequestFormData(method, uri string, postData map[string]string) (resp *HttpResponse) {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	for k, v := range postData {
		w.WriteField(k, v)
	}
	w.Close()
	return Request(method, uri, w.FormDataContentType(), body.String())
}

func Request(method, uri, contentType, payload string) (resp *HttpResponse) {
	return Request2(method, uri, contentType, strings.NewReader(payload))
}

func Request2(method, uri, contentType string, reader io.Reader) (resp *HttpResponse) {
	req, err := http.NewRequest(method, uri, reader)
	if err != nil {
		fmt.Println(err)
	}
	if len(contentType) > 0 {
		req.Header.Add("Content-Type", contentType)
	}
	r, err1, header, statusCode := do(req)
	resp = &HttpResponse{}
	resp.err = err1
	resp.header = header
	resp.result = r
	resp.StatusCode = statusCode
	return resp
}

func GET(uri string, params url.Values) (resp *HttpResponse) {
	if params != nil {
		uri += "?" + params.Encode()
	}
	return Request("GET", uri, "", "")
}

func DELETE(uri string, params url.Values) (resp *HttpResponse) {
	return Request("DELETE", uri, "", params.Encode())
}

func do(req *http.Request) (result []byte, err error, header http.Header, StatusCode int) {
	tr := http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	if os.Getenv("http_proxy") != "" {
		proxyUrl, err := url.Parse(os.Getenv("http_proxy"))
		if err != nil {
			fmt.Println(err)
		}
		tr.Proxy = http.ProxyURL(proxyUrl)
	} else if os.Getenv("HTTP_PROXY") != "" {
		proxyUrl, err := url.Parse(os.Getenv("HTTP_PROXY"))
		if err != nil {
			fmt.Println(err)
		}
		tr.Proxy = http.ProxyURL(proxyUrl)
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse /* 不进入重定向 */
		},
		Transport: &tr,
		Timeout:   600 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err, nil, -1
	}
	if res.StatusCode == 204 {
		return nil, nil, res.Header, res.StatusCode
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if res.StatusCode > 400 && err == nil {
		err = errors.New(string(body))
	}
	return body, err, res.Header, res.StatusCode
}

func makePayload(params interface{}) (string, error) {
	var payload string
	paramsValue := reflect.Indirect(reflect.ValueOf(params))
	switch paramsValue.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice:
		{
			b, e := json.Marshal(params)
			return string(b), e
		}
	case reflect.String:
		{
			payload = params.(string)
		}
	}
	return payload, nil
}

// 上传文件
// uri                请求地址
// params        post form里数据
// files  key 请求地址上传文件对应field value  上传文件路径
// file               文件
func PostFiles(uri string, params map[string]string, files map[string]string) *HttpResponse {
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)

	for nameField, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			continue
		}
		formFile, err := writer.CreateFormFile(nameField, file.Name())
		if err != nil {
			return &HttpResponse{err: err}
		}

		_, err = io.Copy(formFile, file)
		if err != nil {
			return &HttpResponse{err: err}
		}
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err := writer.Close()
	if err != nil {
		return &HttpResponse{err: err}
	}
	return Request("POST", uri, writer.FormDataContentType(), body.String())
}

// 上传文件
// uri                请求地址
// params        post form里数据
// files  key 请求地址上传文件对应field value  上传文件路径
// file               文件

func PostFiles3(uri string, params map[string]string, files map[string]io.Reader) *HttpResponse {
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)

	for nameField, file := range files {

		formFile, err := writer.CreateFormFile(nameField, time.Now().String())
		if err != nil {
			return &HttpResponse{err: err}
		}

		_, err = io.Copy(formFile, file)
		if err != nil {
			return &HttpResponse{err: err}
		}
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err := writer.Close()
	if err != nil {
		return &HttpResponse{err: err}
	}
	return Request("POST", uri, writer.FormDataContentType(), body.String())
}

type File struct {
	Reader io.Reader
	Name   string
}

// 上传文件
// uri                请求地址
// params        post form里数据
// files  key 请求地址上传文件对应field value  上传文件路径
// file               文件
func PostFiles2(uri string, params map[string]string, files map[string]*File) *HttpResponse {
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)

	for nameField, file := range files {
		formFile, err := writer.CreateFormFile(nameField, file.Name)
		if err != nil {
			return &HttpResponse{err: err}
		}

		_, err = io.Copy(formFile, file.Reader)
		if err != nil {
			return &HttpResponse{err: err}
		}
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err := writer.Close()
	if err != nil {
		return &HttpResponse{err: err}
	}
	return Request("POST", uri, writer.FormDataContentType(), body.String())
}

// 上传文件
// uri                请求地址
// params        post form里数据
// files  key 请求地址上传文件对应field value  上传文件路径
// file               文件
func PutFile(uri string, file *os.File) *HttpResponse {
	//body := new(bytes.Buffer)
	body, err := ioutil.ReadAll(file)
	if err != nil {
		return &HttpResponse{err: err}
	}

	return Request2("PUT", uri, "application/octet-stream", bytes.NewReader(body))
}
