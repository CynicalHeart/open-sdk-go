package constant

type ResultCode struct {
	status string
	msg    string
}

func (r ResultCode) Status() string {
	return r.status
}

func (r ResultCode) Msg() string {
	return r.msg
}

var (
	SUCCESS                   = ResultCode{status: "M0200", msg: "操作成功"}
	TOO_MANY_REQUESTS         = ResultCode{status: "M0429", msg: "请求次数过多，请稍后重试"}
	UNSUPPORTED_REQUEST_TYPE  = ResultCode{status: "M0430", msg: "不支持的请求类型"}
	FAILED                    = ResultCode{status: "M0500", msg: "系统繁忙，请稍后重试"}
	VALIDATE_FAILED           = ResultCode{status: "M0555", msg: "参数校验失败"}
	CALL_FAILED               = ResultCode{status: "M0511", msg: "三方服务调用失败"}
	GET_INTERFACE_INFO_FAILED = ResultCode{status: "M0512", msg: "获取接口信息失败"}
	RESPONSE_CONVERSION_ERROR = ResultCode{status: "M0513", msg: "响应信息转换失败"}
	REQUEST_PARAM_NOT_NULL    = ResultCode{status: "M0514", msg: "接口请求参数不能为空"}
	REQUEST_HEADER_NOT_NULL   = ResultCode{status: "M0515", msg: "接口请求头不能为空"}
)
