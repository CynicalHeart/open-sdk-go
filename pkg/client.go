package opensdk

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	cons "open-sdk-go/internal/constant"
	"open-sdk-go/internal/encrypt"
)

type AlgorithmType string

const (
	RSA AlgorithmType = "RSA"
	SM2 AlgorithmType = "SM2"
)

var logger = log.New(os.Stdout, "【OpenSDK】", log.Lshortfile|log.Ldate|log.Ltime)

// 请求客户端
type OpenPlatformClient[T any] struct {
	AppKey        string            `json:"appKey"`
	AppSecret     string            `json:"appSecret"`
	RequestUrl    string            `json:"requestUrl"`
	RequestData   T                 `json:"requestData,omitempty"`
	Headers       map[string]string `json:"header,omitempty"`
	AlgorithmType AlgorithmType     `json:"algorithmType"`
	IsReport      bool              `json:"isReport"`
	ProductCase   string            `json:"productCase"`
}

func NewClient[T any](appKey, appSecret, requestUrl string) *OpenPlatformClient[T] {
	return &OpenPlatformClient[T]{
		AppKey:        appKey,
		AppSecret:     appSecret,
		RequestUrl:    requestUrl,
		AlgorithmType: RSA, // 默认算法类型为RSA
		IsReport:      false,
	}
}

func (c *OpenPlatformClient[T]) SetData(data T) *OpenPlatformClient[T] {
	c.RequestData = data
	return c
}

func (c *OpenPlatformClient[T]) SetHeader(header map[string]string) *OpenPlatformClient[T] {
	if c.Headers == nil {
		c.Headers = make(map[string]string)
	}
	c.Headers = header
	return c
}

func (c *OpenPlatformClient[T]) SetAlgorithmType(algorithmType AlgorithmType) *OpenPlatformClient[T] {
	c.AlgorithmType = algorithmType
	return c
}

func (c *OpenPlatformClient[T]) SetReport(isReport bool) *OpenPlatformClient[T] {
	c.IsReport = isReport
	return c
}

func (c *OpenPlatformClient[T]) SetProductCase(productCase string) *OpenPlatformClient[T] {
	c.ProductCase = productCase
	return c
}

func (c *OpenPlatformClient[T]) String() string {
	return fmt.Sprintf("OpenPlatformClient{AppKey:%s, AppSecret:%s, RequestUrl:%s, RequestData:%v, Headers:%v, AlgorithmType:%s, IsReport:%v, ProductCase:%s}",
		c.AppKey, c.AppSecret, c.RequestUrl, c.RequestData, c.Headers, c.AlgorithmType, c.IsReport, c.ProductCase)
}

// 向云枢发起请求
func (c *OpenPlatformClient[T]) Send() *BaseResult[any] {
	// 切割请求地址，获取域名[:port]
	if c.RequestUrl == "" {
		logger.Println("请求地址不能为空")
		return failWithCodeAndMsg[any]("M0514", "请求地址不能为空")
	}
	_, err := url.ParseRequestURI(c.RequestUrl)
	if err != nil {
		logger.Println("请求地址格式不正确")
		return failWithCodeAndMsg[any]("M0514", "请求地址格式不正确")
	}
	reg, _ := regexp.Compile("/")
	idx := reg.FindAllStringIndex(c.RequestUrl, 3)
	url := c.RequestUrl[:idx[2][0]]
	if strings.HasPrefix(url, "https://open.yljr.com") {
		url = url + "/api"
	}
	url = url + "/api-app/sdk/request"

	// 序列化请求体
	jsonStr, _ := json.Marshal(c)
	enc, err := encrypt.RsaEncode(string(jsonStr))
	if err != nil {
		logger.Printf("报文加密失败, 失败原因: %s.", err.Error())
		return failWithCodeAndMsg[any]("M0514", err.Error())
	}
	requestMap := map[string]string{
		"encryptData": enc,
	}
	requestBody, _ := json.Marshal(requestMap)

	// 加签
	sign, err := c.secure()
	if err != nil {
		logger.Printf("报文加签失败, 失败原因: %s.", err.Error())
		return failWithCodeAndMsg[any]("M0514", err.Error())
	}
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("algorithm", string(c.AlgorithmType))
	request.Header.Set("secure", sign)

	// 添加自定义header
	for k, v := range c.Headers {
		request.Header.Set(k, v)
	}

	// 发送请求
	client := &http.Client{Timeout: 5 * time.Second, Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	response, err := client.Do(request) // 发送请求
	if err != nil {
		logger.Printf("请求失败, 失败原因: %s.", err.Error())
		return failWithResultCode[any](cons.CALL_FAILED)
	}
	defer response.Body.Close()

	// 处理响应并反序列化
	body, _ := io.ReadAll(response.Body)
	var result BaseResult[any]
	if err = json.Unmarshal(body, &result); err != nil {
		logger.Printf("反序列化响应异常, 失败原因: %s.", err.Error())
		return failWithResultCode[any](cons.RESPONSE_CONVERSION_ERROR)
	}
	logger.Printf("请求响应结果:%+v", &result)
	return &result
}

// 对报文加签
func (c *OpenPlatformClient[T]) secure() (string, error) {

	var sign string
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	headerMap := map[string]interface{}{
		"appKey":    c.AppKey,
		"appSecret": c.AppSecret,
		"timestamp": timestamp,
	}

	headerJson, _ := json.Marshal(headerMap)
	if c.AlgorithmType == RSA {
		sign, _ = encrypt.RsaEncode(string(headerJson))
	} else {
		return sign, fmt.Errorf("当前语言不支持的算法类型")
	}

	return sign, nil
}
