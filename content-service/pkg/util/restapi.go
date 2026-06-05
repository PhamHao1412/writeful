package util

import (
	"encoding/json"
	"net/http"

	"log/slog"

	"github.com/sendgrid/rest"
)

func SendRequest(header map[string]string, reqBody interface{}, queryParams map[string]string, baseURL string, method rest.Method) (res []byte, err error, statusCode int) {
	var bdBytes []byte
	if reqBody != nil {
		switch bd := reqBody.(type) {
		case []byte:
			bdBytes = bd
		default:
			bdBytes, err = json.Marshal(reqBody)
			if err != nil {
				return res, err, http.StatusInternalServerError
			}
		}
	}

	request := rest.Request{
		Method:      method,
		BaseURL:     baseURL,
		Body:        bdBytes,
		Headers:     header,
		QueryParams: queryParams,
	}
	response, err := rest.Send(request)
	slog.Info("Send Request", "headers", header)
	if err != nil {
		return res, err, http.StatusInternalServerError
	}

	res = []byte(response.Body)
	statusCode = response.StatusCode
	return res, nil, statusCode
}
