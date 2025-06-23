package showdoc

import (
	http2 "alicode.mukj.cn/yjkj.ink/work/http"
	"fmt"
)


type ServiceDoc struct {
	ErrorCode int `json:"error_code"`
	Data      struct {
		ItemId        string      `json:"item_id"`
		ItemDomain    string      `json:"item_domain"`
		IsArchived    string      `json:"is_archived"`
		ItemName      string      `json:"item_name"`
		DefaultPageId string      `json:"default_page_id"`
		DefaultCatId2 int         `json:"default_cat_id2"`
		DefaultCatId3 int         `json:"default_cat_id3"`
		DefaultCatId4 interface{} `json:"default_cat_id4"`
		UnreadCount   interface{} `json:"unread_count"`
		ItemType      string      `json:"item_type"`
		Menu          struct {
			Pages    []interface{} `json:"pages"`
			Catalogs []struct {
				CatId       string `json:"cat_id"`
				CatName     string `json:"cat_name"`
				ItemId      string `json:"item_id"`
				SNumber     string `json:"s_number"`
				Addtime     string `json:"addtime"`
				ParentCatId string `json:"parent_cat_id"`
				Level       string `json:"level"`
				Pages       []struct {
					PageId    string `json:"page_id"`
					AuthorUid string `json:"author_uid"`
					CatId     string `json:"cat_id"`
					PageTitle string `json:"page_title"`
					Addtime   string `json:"addtime"`
					ExtInfo   string `json:"ext_info"`
				} `json:"pages"`
				Catalogs []struct {
					CatId       string `json:"cat_id"`
					CatName     string `json:"cat_name"`
					ItemId      string `json:"item_id"`
					SNumber     string `json:"s_number"`
					Addtime     string `json:"addtime"`
					ParentCatId string `json:"parent_cat_id"`
					Level       string `json:"level"`
					Pages       []struct {
						PageId    string `json:"page_id"`
						AuthorUid string `json:"author_uid"`
						CatId     string `json:"cat_id"`
						PageTitle string `json:"page_title"`
						Addtime   string `json:"addtime"`
						ExtInfo   string `json:"ext_info"`
					} `json:"pages"`
					Catalogs []interface{} `json:"catalogs"`
				} `json:"catalogs"`
			} `json:"catalogs"`
		} `json:"menu"`
		IsLogin       bool          `json:"is_login"`
		ItemEdit      bool          `json:"item_edit"`
		ItemManage    bool          `json:"item_manage"`
		ItemPermn     bool          `json:"ItemPermn"`
		ItemCreator   bool          `json:"ItemCreator"`
		CurrentPageId string        `json:"current_page_id"`
		GlobalParam   []interface{} `json:"global_param"`
		ShowWatermark string        `json:"show_watermark"`
	} `json:"data"`
}

func (doc *ShowDoc) RefreshServiceInfo(itemId string) *ServiceDoc {
	if doc.Header == nil {
		return nil
	}
	item := map[string]string{"item_id": itemId,"user_token":""}
	resp := http2.POSTFormDataWithHeader(fmt.Sprintf("%s%s", doc.Host, infoMethod), item, &doc.Header)
	if resp.Error() != nil {
		fmt.Println("获取文档详情失败1,", resp.Error())
		return nil
	}
	var result *ServiceDoc
	err := resp.Resp(&result)
	if err != nil {
		fmt.Println("获取文档详情失败2,", string(resp.Byte()))
		return nil
	}
	if result.ErrorCode == 0 {
		//删除原文件夹
		for _, cat := range result.Data.Menu.Catalogs {
			err = doc.DeleteCat(itemId, cat.CatId)
			if err != nil {
				fmt.Println("删除服务文档失败：", err)
			}
		}
		return result
	}
	fmt.Println("获取文档详情失败3,", string(resp.Byte()))
	return nil
}

func (doc *ShowDoc) DeleteCat(itemId,cat_id string) error {
	if doc.Header == nil {
		return nil
	}
	item := map[string]string{"item_id": itemId,"user_token":"","cat_id":cat_id}
	resp := http2.POSTFormDataWithHeader(fmt.Sprintf("%s%s", doc.Host, deleteMethod), item, &doc.Header)
	if resp.Error() != nil {
		fmt.Println("删除项目文档失败1,", resp.Error())
		return resp.Error()
	}
	return nil
}