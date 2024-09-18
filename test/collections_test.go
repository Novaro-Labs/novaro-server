package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"novaro-server/model"
	"testing"
)

func TestCollectionsTweet(t *testing.T) {
	// 创建一个请求

	collections := model.Collections{
		UserId: "1",
	}
	marshal, _ := json.Marshal(collections)

	req, err := http.NewRequest("POST", "http://localhost:8080/v1/api/collections/add", bytes.NewBuffer(marshal))
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

func TestCollectionsTweetRm(t *testing.T) {
	// 创建一个请求

	//collections := model.Collections{
	//	UserId: "1",
	//	PostId: "1",
	//}
	//marshal, _ := json.Marshal(collections)
	//
	//req, err := http.NewRequest("POST", "http://localhost:8080/v1/api/collections/remove", bytes.NewBuffer(marshal))
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//// 设置请求头
	//req.Header.Set("Content-Type", "application/json")
	//
	//// 创建 HTTP 客户端
	//client := &http.Client{}
	//
	//// 发送请求
	//resp, err := client.Do(req)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//defer resp.Body.Close()

}
