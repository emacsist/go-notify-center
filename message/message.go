package message

// Message : common
type Message struct {
	MessageID     string   `json:"messageID"`
	From          string   `json:"from"`
	To            []string `json:"to"`
	Body          string   `json:"body"`
	CallbackQueue string   `json:"callbackQueue"`
}

// CallbackData : 回调数据
type CallbackData struct {
	MessageID     string   `json:"messageID"`
	ErrorCode     int      `json:"error_code"`
	ErrorInfo     string   `json:"error_info"`
	CallbackQueue string   `json:"callbackQueue"`
	ToOK          []string `json:"toOK"`
	ToError       []string `json:"toError"`
	From          string   `json:"from"`
}

// Email : 邮箱
type Email struct {
	Message
	Subject string `json:"subject"`
}

// BuildCallbackData : 构建回调对象
func BuildCallbackData(messageID string, errorCode int, errorInfo string, callbackQueue string, toOK []string, toError []string, from string) CallbackData {
	var result CallbackData
	result.CallbackQueue = callbackQueue
	result.ErrorCode = errorCode
	result.ErrorInfo = errorInfo
	result.MessageID = messageID
	result.ToOK = toOK
	result.ToError = toError
	result.From = from
	return result
}
