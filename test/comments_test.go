package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"novaro-server/model"
	"testing"
)

func TestAddCommnets(t *testing.T) {
	// 创建一个请求

	body := model.Comments{
		UserId:  "1",
		PostId:  "1",
		Content: "xxxxxxxxxx22222xxxxxxxxxxxxxxx",
	}
	marshal, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", "http://localhost:8080/v1/api/comments/add", bytes.NewBuffer(marshal))
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
	defer resp.Body.Close()
}

func TestGetCommentByPostId(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8080/v1/api/comments/getCommentsListByPostId?postId=1", nil)
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

func TestGetCommentsListByParentId(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8080/v1/api/comments/getCommentsListByParentId?parentId=1", nil)
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
func TestGetCommentsListByUserId(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8080/v1/api/comments/getCommentsListByUserId?userId=1", nil)
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

func TestDelete(t *testing.T) {
	req, err := http.NewRequest("DELETE", "http://localhost:8080/v1/api/comments/delete?id=", nil)
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
