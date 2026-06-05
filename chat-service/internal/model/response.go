package model

type TResponse[T any] struct {
	Data T `json:"data"`
}
