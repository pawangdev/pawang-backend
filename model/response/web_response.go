package response

type WebResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResponseSuccess(success bool, message string, data interface{}) WebResponse {
	response := WebResponse{
		Success: success,
		Message: message,
		Data:    data,
	}

	return response
}
