package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"novaro-server/model"
	"testing"
)

func TestCollectionsTweet(t *testing.T) {
	// 创建一个请求

	collections := model.Collections{
		UserId: "1",
		PostId: "1",
	}
	marshal, _ := json.Marshal(collections)

	req, err := http.NewRequest("POST", "/v1/api/collections/add", bytes.NewBuffer(marshal))
	if err != nil {
		t.Fatal(err)
	}

	// 创建一个 ResponseRecorder (which satisfies http.ResponseWriter) 来记录响应
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// 你的测试逻辑

	})

	// 直接调用服务端函数，并传入 ResponseRecorder 和 Request
	handler.ServeHTTP(rr, req)

	// 检查返回的状态码是否是期望的 200
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// 检查返回的响应体是否是期望的结果
	expected := `{"message":"your expected response"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
