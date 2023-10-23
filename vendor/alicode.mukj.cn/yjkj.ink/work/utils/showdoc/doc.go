package showdoc

import (
	"alicode.mukj.cn/yjkj.ink/work/http"
	"alicode.mukj.cn/yjkj.ink/work/utils"
	"errors"
	"fmt"
)

type ApiKey struct {
	ApiKey   string `json:"api_key"`
	ApiToken string `json:"api_token"`
}

type API struct {
	CatName     string `json:"cat_name"`
	PageTitle   string `json:"page_title"`
	PageContent string `json:"page_content"`
}

type APIDoc struct {
	ApiKey      string `json:"api_key"`
	ApiToken    string `json:"api_token"`
	CatName     string `json:"cat_name"`
	PageTitle   string `json:"page_title"`
	PageContent string `json:"page_content"`
}

func (doc *ShowDoc) WriteToApiMarkDown(api *API) error {
	apiDoc := &APIDoc{}
	utils.CopyTo(api, apiDoc)
	if doc.ApiKey == nil {
		return errors.New("showdoc 上传接口文档error")
	}
	apiDoc.ApiKey = doc.ApiKey.ApiKey
	apiDoc.ApiToken = doc.ApiKey.ApiToken
	uri := fmt.Sprintf("%s/server/index.php?s=/api/item/updateByApi", doc.Host)
	resp := http.POSTJson(uri, apiDoc)
	if resp.Error() != nil {
		return resp.Error()
	}
	return nil
}

func (doc *ShowDoc) GetServiceByName(name string) *ItemData {
	items := doc.GetServiceList()
	if items == nil {
		return nil
	}
	for _, item := range items.Data {
		if item.ItemName == name {
			return item
		}
	}
	return nil
}

func (doc *ShowDoc) CreateApiKey(name string) *ApiKey {
	if doc.ApiKey != nil {
		return doc.ApiKey
	}
	{
		item := doc.GetServiceByName(name)
		if item != nil {

			err := doc.DeleteService(item.ItemID)
			if err != nil {
				fmt.Println("删除服务文档失败：", err)
				resp := doc.GetApiKey(item.ItemID)
				if resp != nil {
					apiKey := &ApiKey{ApiKey: resp.Data.APIKey, ApiToken: resp.Data.APIToken}
					doc.ApiKey = apiKey
					return apiKey
				}
				return nil
			}
		}
	}

	err := doc.AddService(name)
	if err != nil {
		fmt.Println("添加服务文档失败：", err)
		return nil
	}
	item := doc.GetServiceByName(name)
	if item != nil {
		resp := doc.GetApiKey(item.ItemID)
		if resp != nil {
			apiKey := &ApiKey{ApiKey: resp.Data.APIKey, ApiToken: resp.Data.APIToken}
			doc.ApiKey = apiKey
			return apiKey
		}
	}
	return nil
}
