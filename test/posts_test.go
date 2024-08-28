package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"novaro-server/model"
	"testing"
)

func TestGetPostsById(t *testing.T) {

	req, err := http.NewRequest("GET", "http://localhost:8080/v1/api/posts/getPostsById?id=1", nil)
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

func TestGetPostsByUserId(t *testing.T) {

	req, err := http.NewRequest("GET", "http://localhost:8080/v1/api/posts/getPostsByUserId?userId=1", nil)
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

func TestGetPostsList(t *testing.T) {
	req, err := http.NewRequest("POST", "http://localhost:8080/v1/api/posts/getPostsList", nil)
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

func TestSavePosts(t *testing.T) {
	c := model.Posts{
		UserId:  "1",
		Content: "xxxxxxx11111111111111111111111",
	}

	marshal, _ := json.Marshal(c)

	req, err := http.NewRequest("POST", "http://localhost:8080/v1/api/posts/savePosts", bytes.NewBuffer(marshal))
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

func TestDelPostsById(t *testing.T) {

	req, err := http.NewRequest("DELETE", "http://localhost:8080/v1/api/posts/delPostsById?id=1", nil)
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
