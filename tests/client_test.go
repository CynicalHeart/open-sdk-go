package tests

import (
	opensdk "open-sdk-go/pkg"
	"testing"
)

func TestClient(t *testing.T) {

	var appKey string = "xxx"     // 云枢 - 个人中心获取
	var appSecret string = "xxx"  // 云枢 - 个人中心获取
	var requestUrl string = "xxx" // 云枢 - 申请沙箱、申请测试、申请上线获取

	// 1. 创建客户端
	client := opensdk.NewClient[any](appKey, appSecret, requestUrl)
	// 2. 设置数据：Data可为任意类型（能序列化对象），如：string、map[string]string、struct{}等
	client.SetData(map[string]string{
		"param1": "value1",
		"param2": "value2",
		"param3": "value3",
	})
	// 3. 自定义请求头(若有)
	client.SetHeader(map[string]string{})
	// 4. 发送请求
	result := client.Send()
	// 5. 处理结果，Data为map[string]interface{}类型
	t.Logf("result: %+v, type of data: %T", result, result.Data)
}
