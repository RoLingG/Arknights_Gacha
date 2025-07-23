package utils

import (
	"Go_Arknights_Gacha/global"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// GachaHistoryParams 定义查询参数的结构体
type GachaHistoryParams struct {
	UID      string `url:"uid"`
	Category string `url:"category"`
	Size     string `url:"size"`
	GachaTs  string `url:"gachaTs, omitempty"`
}

// GachaHistoryResp 定义了整个响应的结构
type GachaHistoryResponse struct {
	Code int              `json:"code"`
	Data GachaHistoryData `json:"data"`
	Msg  string           `json:"msg"`
}

// Data 定义了响应中的数据部分
type GachaHistoryData struct {
	List    []CharInfo `json:"list"`
	HasMore bool       `json:"hasMore"`
}

// CharInfo 定义了角色信息的结构
type CharInfo struct {
	PoolID   string `json:"poolId"`
	PoolName string `json:"poolName"`
	CharID   string `json:"charId"`
	CharName string `json:"charName"`
	Rarity   int    `json:"rarity"`
	IsNew    bool   `json:"isNew"`
	GachaTs  string `json:"gachaTs"`
	Pos      int    `json:"pos"`
}

func GachaHistoryGet(akCookie *http.Cookie, uid, category, token, gachaTs, pos, u8Token string) GachaHistoryData {
	params := GachaHistoryParams{
		UID:      uid,
		Category: category,
		Size:     "10",
	}

	// 将结构体编码为查询字符串
	query := url.Values{}
	query.Set("uid", params.UID)
	query.Set("category", params.Category)
	query.Set("size", params.Size)
	if gachaTs != "" && pos != "" {
		// 只有当 GachaTs 和 pos 有值时才添加到查询字符串中
		query.Set("gachaTs", gachaTs)
		query.Set("pos", pos)
	}

	// 完整 URL 包含查询参数
	fullURL := global.GachaHistory_URL + "?" + query.Encode()

	// 创建一个新的 HTTP 客户端
	client := &http.Client{}

	// 创建一个新的 GET 请求
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return GachaHistoryData{}
	}

	// 设置请求 Header
	req.Header.Set("x-account-token", token)
	req.Header.Set("x-role-token", u8Token)
	req.AddCookie(akCookie)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return GachaHistoryData{}
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to fetch data:", resp.Status)
		return GachaHistoryData{}
	}

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return GachaHistoryData{}
	}

	// 解码 JSON 响应体
	var response GachaHistoryResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return GachaHistoryData{}
	}
	data := response.Data
	return data
}
