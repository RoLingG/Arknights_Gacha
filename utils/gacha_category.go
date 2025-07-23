package utils

import (
	"Go_Arknights_Gacha/global"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type CategoryParams struct {
	Uid string `url:"uid"`
}

type GachaCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CategoryResp struct {
	Data []GachaCategory `json:"data"`
	Msg  string          `json:"msg"`
	Code int             `json:"code"`
}

type Category struct {
	ID string `json:"id"`
}

func GachaCategoryGet(akCookie *http.Cookie, uid, u8Token, token string) []Category {
	// 创建查询参数实例并填充数据
	params := CategoryParams{
		Uid: uid,
	}

	query := url.Values{}
	query.Add("uid", params.Uid)

	// 完整 URL 包含查询参数
	fullURL := global.GachaCate_URL + "?" + query.Encode()

	// 创建一个新的 HTTP 客户端
	client := &http.Client{}

	// 创建一个新的 GET 请求
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}
	req.AddCookie(akCookie)

	// 设置请求 Header
	req.Header.Add("X-Account-Token", token)
	req.Header.Add("X-Role-Token", u8Token)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to fetch data: %s\n", resp.Status)
		return nil
	}

	// 检查 Content-Type
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		fmt.Println("Unexpected content type:", contentType)
		return nil
	}

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	// 检查响应头是否为 gzip 压缩
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(bytes.NewReader(body))
		if err != nil {
			fmt.Println("Error decompressing response body:", err)
			return nil
		}
		defer reader.Close()
		body, err = ioutil.ReadAll(reader)
		if err != nil {
			fmt.Println("Error reading decompressed response body:", err)
			return nil
		}
	}

	// 将 JSON 字符串解码到一个结构体中
	var responseData CategoryResp
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Println("Error unMarshaling JSON:", err)
		return nil
	}

	// 检查 API 返回的错误信息
	if responseData.Code != 0 {
		fmt.Printf("API returned error: %s\n", responseData.Msg)
		return nil
	}

	// 构造返回的 Category 切片
	var cate []Category
	for _, v := range responseData.Data {
		cate = append(cate, Category{ID: v.ID})
	}
	return cate
}
