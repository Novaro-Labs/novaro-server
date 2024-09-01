package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"novaro-server/api"
	"testing"
)

func TestTagRecords(t *testing.T) {

	recordsApi := api.TagsRecordsApi{
		TagId:  "1",
		PostId: "1",
	}

	marshal, _ := json.Marshal(recordsApi)

	req, err := http.NewRequest("POST", "http://localhost:8080/v1/api/tags/records/add", bytes.NewBuffer(marshal))
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
