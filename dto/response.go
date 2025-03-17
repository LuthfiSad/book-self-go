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

type PaginatedResponseData[T any] struct {
	Timestamp  time.Time `json:"timestamp"`
	Message    any       `json:"message,omitempty"`
	Error      any       `json:"error,omitempty"`
	Data       T         `json:"data"`
	Page       int       `json:"page,omitempty"`
	PerPage    int       `json:"per_page,omitempty"`
	TotalPages int       `json:"total_pages,omitempty"`
	TotalItems int64     `json:"total_items,omitempty"`
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

func NewPaginatedResponseData[T any](data T, page, perPage, totalPages int, totalItems int64) PaginatedResponseData[T] {
	return PaginatedResponseData[T]{
		Timestamp:  time.Now(),
		Message:    "success",
		Data:       data,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
		TotalItems: totalItems,
	}
}
