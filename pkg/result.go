package opensdk

import (
	"fmt"
	cons "open-sdk-go/internal/constant"
)

type BaseResult[T any] struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Data   T      `json:"data,omitempty"`
}

// 创建基本结果
func newBaseResult[T any](status string, msg string) *BaseResult[T] {
	return &BaseResult[T]{Status: status, Msg: msg}
}

// 创建成功结果
func success[T any]() *BaseResult[T] {
	return newBaseResult[T](cons.SUCCESS.Status(), cons.SUCCESS.Msg())
}

// 创建成功结果并包含数据
func successWithData[T any](data T) *BaseResult[T] {
	result := success[T]()
	result.Data = data
	return result
}

// 创建失败结果
func fail[T any]() *BaseResult[T] {
	return newBaseResult[T](cons.FAILED.Status(), cons.FAILED.Msg())
}

// 创建失败结果并包含数据
func failWithData[T any](data T) *BaseResult[T] {
	result := fail[T]()
	result.Data = data
	return result
}

// 使用指定的 ResultCode 创建失败结果
func failWithResultCode[T any](resultCode cons.ResultCode) *BaseResult[T] {
	return newBaseResult[T](resultCode.Status(), resultCode.Msg())
}

// 使用指定的 code 和 msg 创建失败结果
func failWithCodeAndMsg[T any](code, msg string) *BaseResult[T] {
	return newBaseResult[T](code, msg)
}

// 设置数据并返回结果
func (br *BaseResult[T]) setData(data T) *BaseResult[T] {
	br.Data = data
	return br
}

func (br *BaseResult[T]) String() string {
	return fmt.Sprintf(`{"status": "%s", "msg": "%s", "data": %v}`, br.Status, br.Msg, br.Data)
}
