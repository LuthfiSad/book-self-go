package dto

import "time"

type ResponseMessage struct {
	Timestamp time.Time `json:"timestamp"`
	Message   any       `json:"message,omitempty"`
	Error     any       `json:"error,omitempty"`
}

type ResponseData[T any] struct {
	Timestamp time.Time `json:"timestamp"`
	Message   any       `json:"message,omitempty"`
	Error     any       `json:"error,omitempty"`
	Data      T         `json:"data"`
}

func NewResponseMessage(message any) ResponseMessage {
	return ResponseMessage{
		Timestamp: time.Now(),
		Message:   message,
	}
}

func NewResponseData[T any](data T) ResponseData[T] {
	return ResponseData[T]{
		Timestamp: time.Now(),
		Message:   "success",
		Data:      data,
	}
}
