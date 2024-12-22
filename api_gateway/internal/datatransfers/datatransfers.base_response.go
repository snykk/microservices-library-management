package datatransfers

type response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func ResponseSuccess(message string, data interface{}) response {
	return response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ResponseError(message string, err interface{}) response {
	switch e := err.(type) {
	case error:
		err = e.Error()
	default:
		// No action required; keep err as is
	}

	return response{
		Success: false,
		Message: message,
		Errors:  err,
	}
}
