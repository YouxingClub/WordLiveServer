package main

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"net"
	"time"
)

type echoAPI struct {
	Action string `json:"action"`
	Target string `json:"target"`
	From   struct {
		Name      string `json:"name"`
		UUID      string `json:"uuid"`
		Type      string `json:"type"`
		Timestamp int64  `json:"timestamp"`
	} `json:"from"`
	Data interface{} `json:"data,omitempty"`
}

func NewEchoAPI() *echoAPI {
	newUUID := uuid.NewV1()
	ltime := time.Now()
	return &echoAPI{
		Action: "ping",
		Target: newUUID.String(),
		From: struct {
			Name      string `json:"name"`
			UUID      string `json:"uuid"`
			Type      string `json:"type"`
			Timestamp int64  `json:"timestamp"`
		}{
			Name:      newUUID.String(),
			UUID:      newUUID.String(),
			Type:      "server",
			Timestamp: ltime.UnixMilli(),
		},
		Data: nil,
	}
}

// 返回echoAPI类型
func (e *echoAPI) GetStruct() echoAPI {
	return *e
}

// json序列化
func (e *echoAPI) Marshal() ([]byte, error) {
	// 首先使用 json.Marshal 将结构体序列化
	jsonData, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	// 将序列化后的 JSON 转换为 map 以便我们可以插入空对象
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, err
	}

	// 如果 Data 字段被 omitempty 删除，手动添加空对象
	if _, exists := result["data"]; !exists {
		result["data"] = map[string]interface{}{}
		// 删除target
		delete(result, "target")
	}

	// 再次序列化以生成最终的 JSON 字符串
	finalJSON, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return finalJSON, nil
}

// json反序列化
func (e *echoAPI) Unmarshal(data []byte) error {
	return json.Unmarshal(data, e)
}

type Device struct {
	Name      string
	UUID      string
	Type      string
	Timestamp int64
	Addr      net.Addr
}
