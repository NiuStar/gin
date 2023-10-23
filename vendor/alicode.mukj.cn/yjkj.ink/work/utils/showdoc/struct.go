package showdoc

type ItemData struct {
	Creator         int64       `json:"creator"`
	IsDel           string      `json:"is_del"`
	IsPrivate       int64       `json:"is_private"`
	ItemDescription string      `json:"item_description"`
	ItemDomain      string      `json:"item_domain"`
	ItemID          string      `json:"item_id"`
	ItemName        string      `json:"item_name"`
	ItemType        string      `json:"item_type"`
	LastUpdateTime  string      `json:"last_update_time"`
	SNumber         interface{} `json:"s_number"`
	UID             string      `json:"uid"`
}
type ItemResponse struct {
	Data      []*ItemData `json:"data"`
	ErrorCode int64       `json:"error_code"`
}

type ApiKeyResponse struct {
	Data struct {
		Addtime       string `json:"addtime"`
		APIKey        string `json:"api_key"`
		APIToken      string `json:"api_token"`
		ID            string `json:"id"`
		ItemID        string `json:"item_id"`
		LastCheckTime string `json:"last_check_time"`
	} `json:"data"`
	ErrorCode int64 `json:"error_code"`
}
