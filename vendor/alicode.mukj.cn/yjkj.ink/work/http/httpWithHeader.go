package http

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func POSTDataWithHeader(uri string, params url.Values, header2 *http.Header) (resp *HttpResponse) {
	return RequestWithHeader("POST", uri, "", params.Encode(), header2)
}

func PATCHDataWithHeader(uri string, params url.Values, header2 *http.Header) (resp *HttpResponse) {
	return RequestWithHeader("PATCH", uri, "", params.Encode(), header2)
}

func PATCHJsonWithHeader(uri string, params interface{}, header2 *http.Header) (resp *HttpResponse) {
	return RequestJsonBodyWithHeader("PATCH", uri, params, header2)
}

func POSTJsonWithHeader(uri string, params interface{}, header2 *http.Header) (resp *HttpResponse) {
	return RequestJsonBodyWithHeader("POST", uri, params, header2)
}

func RequestJsonBodyWithHeader(method, uri string, params interface{}, header2 *http.Header) (resp *HttpResponse) {
	payload, err := makePayload(params)
	if err != nil {
		return &HttpResponse{err: err}
	}
	return RequestWithHeader(method, uri, "application/json", payload, header2)
}

func PATCHFormDatanWithHeader(uri string, postData map[string]string, header2 *http.Header) (resp *HttpResponse) {
	return RequestFormDataWithHeader("PATCH", uri, postData, header2)
}

func POSTFormDataWithHeader(uri string, postData map[string]string, header2 *http.Header) (resp *HttpResponse) {
	return RequestFormDataWithHeader("POST", uri, postData, header2)
}

func RequestFormDataWithHeader(method, uri string, postData map[string]string, header2 *http.Header) (resp *HttpResponse) {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	for k, v := range postData {
		w.WriteField(k, v)
	}
	w.Close()
	return RequestWithHeader(method, uri, w.FormDataContentType(), body.String(), header2)
}

func POSTFormUrlDecodeWithHeader(uri string, postData map[string]string,header2 *http.Header) (resp *HttpResponse) {
	value := url.Values{}
	for key, d := range postData {
		value.Add(key, d)
	}
	return RequestWithHeader("POST", uri, "application/x-www-form-urlencoded", value.Encode(),header2)
}

func RequestWithHeader(method, uri, contentType, payload string, header2 *http.Header) (resp *HttpResponse) {
	var reader io.Reader
	if len(payload) > 0 {
		reader = strings.NewReader(payload)
	}
	req, err := http.NewRequest(method, uri, reader)
	if err != nil {
		fmt.Println(err)
	}
	if len(contentType) > 0 {
		req.Header.Add("Content-Type", contentType)
	}

	if header2 != nil {
		for name, values := range *header2 {
			name2 := strings.ToLower(name)

			if name2 == "set-cookie" {
				for _, v := range values {
					req.Header.Add("Cookie", v)
				}
			} else {
				for _, v := range values {
					req.Header.Add(name, v)
				}
			}
		}
	}

	r, err1, header, statu:= do(req)
	resp = &HttpResponse{}
	resp.err = err1
	resp.header = header
	resp.result = r
	resp.StatusCode = statu
	return resp
}

func GETWithHeader(uri string, params url.Values, header2 *http.Header) (resp *HttpResponse) {
	return RequestWithHeader("GET", uri, "", params.Encode(), header2)
}

func DELETEWithHeader(uri string, params url.Values, header2 *http.Header) (resp *HttpResponse) {
	return RequestWithHeader("DELETE", uri, "", params.Encode(), header2)
}

// 上传文件
// uri                请求地址
// params        post form里数据
// files  key 请求地址上传文件对应field value  上传文件路径
// file               文件
func PostFilesWithHeader(uri string, params map[string]string, files map[string]string, header2 *http.Header) (resp *HttpResponse) {
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
	return RequestWithHeader("POST", uri, writer.FormDataContentType(), body.String(), header2)
}

// 上传文件
// uri                请求地址
// params        post form里数据
// files  key 请求地址上传文件对应field value  上传文件路径
// file               文件
func PostFilesWithHeader2(uri string, params map[string]string, files map[string]*os.File, header2 *http.Header) (resp *HttpResponse) {
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)

	for nameField, file := range files {
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
	return RequestWithHeader("POST", uri, writer.FormDataContentType(), body.String(), header2)
}
