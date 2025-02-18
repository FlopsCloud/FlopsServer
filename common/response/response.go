package response

const (
	SuccessCode             = 200
	UnauthorizedCode        = 401
	WrongCaptchaCode        = 451
	UserNotExistCode        = 452
	RecordNotExistCode      = 453
	InvalidRequestParamCode = 455
	InsufficientBalanceCode = 456
	UserBalanceNotExistCode = 457
	ServerErrorCode         = 500
)

type Response struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func OK(data interface{}) Response {
	return Response{Code: SuccessCode, Message: "ok", Data: data}
}

func Error(message string) Response {
	return fail(ServerErrorCode, message)
}

func Fail(code int64, message string) Response {
	return fail(code, message)
}

func fail(code int64, message string) Response {
	if message == "" {
		message = "server error"
	}
	return Response{Code: code, Message: message}
}
