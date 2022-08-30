package exception

type Error struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	ErrorCode int         `json:"error_code"`
	Data      interface{} `json:"data"`
}

func ResponseError(success bool, message string, errorcode int, data interface{}) Error {
	errorMsg := Error{
		Success:   success,
		Message:   message,
		ErrorCode: errorcode,
		Data:      data,
	}

	return errorMsg
}
