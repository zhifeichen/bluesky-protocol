package common

/**
	内部错误未负数
 */
const (
	OK                 = -iota
	CONFIG_ARG_ERROR
	CONFIG_NOT_FOUND
	CONFIG_PARSE_ERROR
	READ_FILE_ERROR
	INER_TYPE_INVALID

	CHAIN_HANDLE_MSG_ERROR
)

var (
	ErrorMap = map[int]string{
		OK:                 "OK",
		READ_FILE_ERROR: 	"文件不存在",
		CONFIG_NOT_FOUND:   "配置文件错误",
		CONFIG_ARG_ERROR:   "请指定配置文件",
		CONFIG_PARSE_ERROR: "解析配置文件错误",
		INER_TYPE_INVALID:  "内部错误:类型错误转换",
		CHAIN_HANDLE_MSG_ERROR: "管道处理消息错误",
	}
)

func GetMsgByCode(code int) string {
	if v, ok := ErrorMap[code]; ok {
		return v
	} else {
		return "未知错误"
	}
}
