package apidoc

import (
	"alicode.mukj.cn/yjkj.ink/work/apidoc/markdown"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strings"
)

type Doc struct {
	Name    string             `json:"name"`
	Methods map[string]*Method `json:"methods"`
}

type Method struct {
	NoteName string `json:"note_name"`
	Export   bool   `json:"export"`
}

//Api文档页面
type ApiPage struct {
	ApiKey      string `json:"api_key" comment:"showdoc上生成的apiKey"`
	ApiToken    string `json:"api_token" comment:"showdoc上生成的apiToken"`
	CatName     string `json:"cat_name" comment:"所属文件夹名称"`
	PageTitle   string `json:"page_title" comment:"PageTitle"`
	PageContent string `json:"page_content" comment:"PageContent"`
}

func NewApiPage(pageTitle, content string) *ApiPage {
	return &ApiPage{PageTitle: pageTitle, PageContent: content}
}

//ApiDoc文件夹
type ApiFolder struct {
	CatName string     `json:"cat_name" comment:"文件夹名称"`
	Pages   []*ApiPage `json:"pages" comment:"当前文件夹下的所有页面"`
}

type API struct {
	ApiUri     string                `json:"api_uri" comment:"showdoc的地址"`
	ApiKey     string                `json:"api_key" comment:"showdoc上生成的apiKey"`
	ApiToken   string                `json:"api_token" comment:"showdoc上生成的apiToken"`
	ApiFolders map[string]*ApiFolder `json:"api_folders" comment:"文件夹"`
	Domain     string                `json:"domain"`
}

//showdoc文档地址，文档项目的apiKey和apiToken，在showdoc里生成
func New(apiUri, apiKey, apiToken, domain string) *API {
	return &API{ApiUri: apiUri + "/server/index.php?s=/api/item/updateByApi", ApiKey: apiKey,
		ApiToken: apiToken, ApiFolders: make(map[string]*ApiFolder), Domain: domain}
}

/*
domain:接口前缀地址
apiPath：接口路由地址
folderName：接口文档文件夹名称
methodName：方法名称
methNameNote：方法备注，一般使用方法中文名
request：请求体结构
response：返回体结构
*/
func (api *API) WriteToApiMarkDown(apiPath, folderName, methodName, methNameNote string, request, response interface{}) (string, error) {
	apiDomain := path.Join(api.Domain, apiPath)
	douri, err := url.Parse(api.Domain)
	if err == nil {
		douri.Path = path.Join(douri.Path, apiPath)
		apiDomain = douri.String()
	}
	apimk := markdown.NewMarkDown(methNameNote, "")
	requestType := reflect.TypeOf(request)
	responseType := reflect.TypeOf(response)
	parms(apimk, apiDomain, methodName, apiPath, methNameNote, requestType, responseType)
	page := NewApiPage(methNameNote, apimk.Content())
	page.CatName = folderName
	page.ApiKey = api.ApiKey
	page.ApiToken = api.ApiToken
	_, err = POSTJson(api.ApiUri, page, nil)
	if err != nil {
		return apiDomain, errors.New("添加接口文档失败，请检查showdoc服务是否开启，可使用docker start showdoc命令开启")
	}
	return apiDomain, nil
}

func parms(mk *markdown.MarkDown, domain, methodName, router, note string, request reflect.Type, response reflect.Type) {
	{
		mk.WriteTitle(3, methodName+"  "+note+"\r\n")
		mk.WriteTitle(4, "http 调用方法\r\n")
		mk.WriteCode("URL:  "+domain+"\r\n", "go")
	}

	{
		mk.WriteContent("\r\n请求参数：\r\n")
		var params [][]string
		{
			var param []string
			param = append(param, "参数名")
			param = append(param, "类型")
			param = append(param, "备注")
			param = append(param, "必填")
			params = append(params, param)
		}

		////fmt.Println("params_list:",params_list)
		{
			//index := 0
			_type := request
			params2 := markdown.ParseParms(_type, 0, false)
			params = append(params, params2...)

			mk.WriteForm(params)
		}
	}

	{
		mk.WriteContent("\r\n返回值：\r\n")
		var params [][]string
		{
			var param []string
			param = append(param, "参数名")
			param = append(param, "类型")
			param = append(param, "备注")
			params = append(params, param)
		}
		if response != nil {
			params = append(params, markdown.ParseParms(response, 0, true)...)
		}

		mk.WriteForm(params)
	}
}

func POSTJson(url string, payload interface{}, headers map[string]string) (string, error) {
	client := &http.Client{}

	var payload_body []byte
	switch payload.(type) {
	case string:
		{

			payload_body = []byte(payload.(string))
		}
	default:
		payload_body, _ = json.Marshal(payload)
	}

	////fmt.Println("payload:",string(payload_body))
	req, err := http.NewRequest("POST", url, strings.NewReader(string(payload_body)))

	if err != nil {
		//fmt.Println(err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	//req.Header.Add("User-Agent", "Xcode")
	//req.Header.Add("Accept", "text/x-xml-plist")

	for name, header := range headers {
		req.Header.Add(name, header)
	}

	res, err := client.Do(req)
	if err != nil {
		//fmt.Println("POST error",err)
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//fmt.Println("POST read error",err)
		return "", err
	}

	////fmt.Println("POSTJson return",string(body))
	return string(body), nil
}
