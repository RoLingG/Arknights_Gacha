package utils

import (
	"Go_Arknights_Gacha/global"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// QueryParams 定义查询参数的结构体
type BindingListParams struct {
	GrantToken string `url:"token"`
	AppCode    string `url:"appCode"`
}

// Data 定义了响应中的数据部分
type Data struct {
	List []AppInfo `json:"list"`
}

// AppInfo 定义了应用信息的结构
type AppInfo struct {
	AppCode     string    `json:"appCode"`
	AppName     string    `json:"appName"`
	BindingList []Binding `json:"bindingList"`
}

// Binding 定义了绑定信息的结构
type Binding struct {
	UID             string `json:"uid"`
	IsOfficial      bool   `json:"isOfficial"`
	IsDefault       bool   `json:"isDefault"`
	ChannelMasterId int    `json:"channelMasterId"`
	ChannelName     string `json:"channelName"`
	NickName        string `json:"nickName"`
	IsDeleted       bool   `json:"isDeleted"`
}

type BindingListResp struct {
	Data   Data   `json:"data"`
	Msg    string `json:"msg"`
	Status int    `json:"status"`
}

func BindingListGet(token string, appCode string) (uid string, err error) {
	// 创建查询参数实例并填充数据
	params := BindingListParams{
		GrantToken: token,
		AppCode:    appCode,
	}

	query := url.Values{}
	query.Add("token", params.GrantToken)
	query.Add("appCode", params.AppCode)

	// 完整 URL 包含查询参数
	fullURL := global.Binding_list_URL + "?" + query.Encode()

	// 创建一个新的 HTTP 客户端
	client := &http.Client{}

	// 创建一个新的 GET 请求
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to fetch data:", resp.Status)
		return "", err
	}

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", err
	}

	// 将 JSON 字符串解码到一个结构体中
	var responseData BindingListResp
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Println("Error unMarshaling JSON:", err)
		return "", err
	}

	uid = responseData.Data.List[0].BindingList[0].UID
	return uid, nil
}
