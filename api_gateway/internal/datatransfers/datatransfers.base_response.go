package datatransfers

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ResponseSuccess(message string, data interface{}) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ResponseError(message string, err error) Response {
	return Response{
		Success: false,
		Message: message,
		Data:    err.Error(),
	}
}
