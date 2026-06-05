package util

import (
	"content-service/internal/model"
	"encoding/json"
)

func ParseResponse[T any](res []byte) (T, error) {
	var response model.TResponse[T]
	if err := json.Unmarshal(res, &response); err != nil {
		var zero T
		return zero, err
	}
	return response.Data, nil
}
