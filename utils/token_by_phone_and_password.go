package utils

import (
	"Go_Arknights_Gacha/global"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"os"
)

type TokenByPhoneAndPasswordReq struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type TokenByPhoneAndPasswordResp struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
	Msg    string `json:"msg"`
	Status int    `json:"status"`
	Type   string `json:"type"`
}

func TokenByPhoneAndPasswordPost() (token string, err error) {
	// 将请求数据编码为 JSON
	godotenv.Load(".env")
	var reqType = TokenByPhoneAndPasswordReq{
		Phone:    os.Getenv("PHONE"),
		Password: os.Getenv("PASSWORD"),
	}
	jsonData, err := json.Marshal(reqType)
	if err != nil {
		fmt.Println("Error marshaling data:", err)
		return "", err
	}
	// 创建一个新的 HTTP 客户端
	client := &http.Client{}

	// 创建一个新的 POST 请求
	req, err := http.NewRequest("POST", global.TokenByPhoneAndPassword_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	// 设置请求 Header
	req.Header.Set("Content-Type", "application/json") // 发送请求

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", err
	}

	// 将 JSON 字符串解码到一个结构体中
	var responseData TokenByPhoneAndPasswordResp
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Println("Error unMarshaling JSON:", err)
		return "", err
	}

	// 提取 token 值
	token = responseData.Data.Token
	return token, nil
}
