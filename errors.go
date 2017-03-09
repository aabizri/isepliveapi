package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// An error to return
type JSONError struct {
	// The error code
	ErrorCode string `json:"code"`

	// If not set, this is automatically populated via the ErrorCode
	Message string `json:"message"`

	// If not set, this is omitted
	Description string `json:"description,omitempty"`

	// HTTP Status Code
	httpStatusCode int

	// Sub-errors
	Errors []JSONError `json:"errors,omitempty"`
}

func (err JSONError) Error() string {
	str := fmt.Sprintf("Error %s: %s", err.ErrorCode, err.Message)
	if err.Description != "" {
		str += fmt.Sprintf(" (%s)", err.Description)
	}
	return str
}

func (jerr JSONError) JSON() ([]byte, error) {
	jresp, err := json.Marshal(jerr)
	return jresp, err
}

func (jerr JSONError) JSONWrite(w http.ResponseWriter) error {
	w.WriteHeader(jerr.httpStatusCode)

	jresp, err := jerr.JSON()
	if err != nil {
		return err
	}

	_, err = bytes.NewBuffer(jresp).WriteTo(w)

	log.Println(jerr.Error())
	return err
}

type errorType struct {
	DefaultMessage string
	httpStatusCode int
}

const (
	InternalError                            string = "I"
	InternalErrorTypeAssertionFailed                = "I101"
	InternalErrorWritingResponseWriterFailed        = "I102"
	InternalErrorReadingRequestFailed               = "I103"
	InternalErrorJSONUnmarshallingFailed            = "I110"
	InternalErrorJSONMarshallingFailed              = "I111"
)

var errorCodeMapper map[string]errorType = map[string]errorType{
	"A":    {"Authorization failed", 400},
	"A10":  {"Authorization header absent", 400},
	"A11":  {"Authorization header invalid", 400},
	"A110": {"Authorization header doesn't have strictly two elements", 400},
	"A111": {"Authorization header's scheme isn't set to \"Bearer\"", 400},
	"A20":  {"Invalid token", 400},
	"I":    {"Internal Server Error: Generic", 500},
	"I101": {"Internal Server Error: Type Assertion Failed", 500},
	"I102": {"Internal Server Error: Failed writing to response writer", 500},
	"I103": {"Internal Server Error: Failed reading request", 500},
	"I110": {"Internal Server Error: JSON Unmarshalling Failed", 500},
	"I111": {"Internal Server Error: JSON Marshalling Failed", 500},
}

func NewError(errorCode, message, description string) JSONError {
	err := JSONError{
		ErrorCode:      errorCode,
		Message:        message,
		Description:    description,
		httpStatusCode: errorCodeMapper[errorCode].httpStatusCode,
	}
	if err.Message == "" {
		err.Message = errorCodeMapper[errorCode].DefaultMessage
	}

	return err
}
