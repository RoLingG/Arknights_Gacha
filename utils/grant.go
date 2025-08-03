package utils

import (
	"Go_Arknights_Gacha/global"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type GrantReq struct {
	AppCode string `json:"appCode"`
	Token   string `json:"token"`
	Type    int    `json:"type"`
}

type GrantResp struct {
	Data struct {
		Token string `json:"token"`
		HgId  string `json:"hgId"`
	} `json:"data"`
	Msg    string `json:"msg"`
	Status int    `json:"status"`
	Type   string `json:"type"`
}

func GrantPost(token string) (grantToken string, hgId string, err error) {
	reqType := GrantReq{
		AppCode: "be36d44aa36bfb5b",
		Token:   token,
		Type:    1,
	}
	// 将请求数据编码为 JSON
	jsonData, err := json.Marshal(reqType)
	if err != nil {
		fmt.Println("Error marshaling data:", err)
		return "", "", err
	}
	// 创建一个新的 HTTP 客户端
	client := &http.Client{}

	// 创建一个新的 POST 请求
	req, err := http.NewRequest("POST", global.Grant_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", "", err
	}

	// 设置请求 Header
	req.Header.Set("Content-Type", "application/json") // 发送请求

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", "", err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", "", err
	}

	// 将 JSON 字符串解码到一个结构体中
	var responseData GrantResp
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Println("Error unMarshaling JSON:", err)
		return "", "", err
	}

	// 提取 token, hgId 值
	grantToken = responseData.Data.Token
	hgId = responseData.Data.HgId
	return grantToken, hgId, nil
}
