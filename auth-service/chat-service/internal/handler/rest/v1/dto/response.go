package dto

import "net/http"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Total   *int64      `json:"total,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func ResponseOK(data interface{}) *Response {
	return &Response{
		Success: true,
		Data:    data,
	}
}

func (r *Response) WithMessage(message string) *Response {
	r.Message = message
	return r
}

func (r *Response) TotalItem(total int64) *Response {
	r.Total = &total
	return r
}

func ResponseError(err error, statusCode int) *Response {
	return &Response{
		Success: false,
		Error:   err.Error(),
		Message: getMessageByStatusCode(statusCode),
	}
}

func ResponseBadRequest(err error) *Response {
	return &Response{
		Success: false,
		Error:   err.Error(),
		Message: "Bad Request",
	}
}

func ResponseUnauthorized(err error) *Response {
	return &Response{
		Success: false,
		Error:   err.Error(),
		Message: "Unauthorized",
	}
}

func ResponseForbidden(err error) *Response {
	return &Response{
		Success: false,
		Error:   err.Error(),
		Message: "Forbidden",
	}
}

func ResponseNotFound(err error) *Response {
	return &Response{
		Success: false,
		Error:   err.Error(),
		Message: "Not Found",
	}
}

func getMessageByStatusCode(statusCode int) string {
	switch statusCode {
	case http.StatusBadRequest:
		return "Bad Request"
	case http.StatusUnauthorized:
		return "Unauthorized"
	case http.StatusForbidden:
		return "Forbidden"
	case http.StatusNotFound:
		return "Not Found"
	case http.StatusInternalServerError:
		return "Internal Server Error"
	default:
		return "Error"
	}
}
