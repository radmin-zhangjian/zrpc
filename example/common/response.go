package common

func GetMsg(key int) string {
	return message[key]
}

func Response(code int, message string, data any) any {
	if data == false || data == "" {
		data = map[string]any{}
	}
	return map[string]any{
		"code": code,
		"msg":  message,
		"data": data,
	}
}
