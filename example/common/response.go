package common

func GetMsg(key int) string {
	return message[key]
}

type ResponseModel struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Response(code int, message string, data any) any {
	if data == false || data == "" {
		data = map[string]any{}
	}
	//return map[string]any{
	//	"code": code,
	//	"msg":  message,
	//	"data": data,
	//}
	return ResponseModel{
		Code: code,
		Msg:  message,
		Data: data,
	}
}
