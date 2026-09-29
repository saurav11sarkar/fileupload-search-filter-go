package utils

import (
	"log"
	"net/http"
)

type AppError struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func NewAppError(status int, code string, message string) *AppError {
	return &AppError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func HanldeError(w http.ResponseWriter, err error) {
	if apperr, ok := err.(*AppError); ok {
		JSONError(w, apperr.Status, apperr.Code, apperr.Message)
		return
	}
	log.Println("Internal Error", err.Error())
	JSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
}
