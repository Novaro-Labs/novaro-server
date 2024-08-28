package test

import (
	"fmt"
	"net/http"
	"testing"
)

func TestList(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8080/v1/api/tags/list", nil)
	if err != nil {
		t.Fatal(err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 创建 HTTP 客户端
	client := &http.Client{}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	// 获取响应信息

	fmt.Println(resp)

	defer resp.Body.Close()

}
