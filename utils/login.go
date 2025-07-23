package utils

import (
	"Go_Arknights_Gacha/global"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type LoginReq struct {
	ShareBy    string `json:"share_by"`
	ShareType  string `json:"share_type"`
	SourceFrom string `json:"source_from"`
	Token      string `json:"token"`
}

type LoginResp struct {
	Data struct{} `json:"data"`
	Msg  string   `json:"msg"`
	Code int      `json:"code"`
}

func LoginPost(u8Token string) *http.Cookie {
	var akCookie *http.Cookie

	reqType := LoginReq{
		ShareType:  "",
		ShareBy:    "",
		SourceFrom: "",
		Token:      u8Token,
	}
	// 将请求数据编码为 JSON
	jsonData, err := json.Marshal(reqType)
	if err != nil {
		fmt.Println("Error marshaling data:", err)
		return akCookie
	}
	// 创建一个新的 HTTP 客户端
	client := &http.Client{}

	// 创建一个新的 POST 请求
	req, err := http.NewRequest("POST", global.Login_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return akCookie
	}

	// 设置请求 Header
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return akCookie
	}
	defer resp.Body.Close()

	// 从响应中提取Cookie
	cookies := resp.Cookies()

	// 如果需要将某个Cookie保存下来，可以这样做
	for _, cookie := range cookies {
		if cookie.Name == "ak-user-center" {
			akCookie = cookie
			// 你可以将这个Cookie保存到文件或内存中，供后续请求使用
		}
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to fetch data: %s\n", resp.Status)
		return akCookie
	}

	// 读取响应体
	var body []byte
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, _ := gzip.NewReader(resp.Body)
		defer reader.Close()
		body, err = ioutil.ReadAll(reader)
	} else {
		body, err = ioutil.ReadAll(resp.Body)
	}
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return akCookie
	}

	// 将 JSON 字符串解码到一个结构体中
	var responseData LoginResp
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Println("Error unMarshaling JSON:", err)
		return akCookie
	}

	// 检查 API 返回的错误信息
	if responseData.Code != 0 {
		fmt.Printf("API returned error: %s\n", responseData.Msg)
		return akCookie
	}
	return akCookie
}
