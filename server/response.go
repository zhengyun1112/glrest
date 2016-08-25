package server

const RESPONSE_CODE_OK = 0
const RESPONSE_CODE_INTERNAL_ERROR = -1
const RESPONSE_MESSAGE_OK = "OK"
const RESPONSE_MESSAGE_INTERNAL_ERROR = "internal server error"

type Response struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	DevMessage string      `json:"dev_message"`
}

func NewOkResponse(data interface{}) *Response {
	return &Response{
		Code:    RESPONSE_CODE_OK,
		Message: RESPONSE_MESSAGE_OK,
		Data:    data,
	}
}

func NewResponse(code int, message string, data interface{}, devMessage string) *Response {
	return &Response{
		Code:       code,
		Message:    message,
		Data:       data,
		DevMessage: devMessage,
	}
}
