package errs

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type ErrorResponse struct {
	RequestId string       `json:"request_id,omitempty"`
	Message   string       `json:"message"`
	Errors    []FieldError `json:"errors,omitempty"`
}

type FieldError struct {
	Field  string `json:"field"`
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

type ErrorBuilder struct {
	StatusCode    int
	Message       string
	Errors        []FieldError
	OriginalError error
	MetaData      interface{}
	RequestId     string
}

func (builder *ErrorBuilder) Error() string {
	if builder.OriginalError != nil {
		return builder.OriginalError.Error()
	}
	return builder.Message
}

func (builder *ErrorBuilder) SetStatusCode(statusCode int) *ErrorBuilder {
	builder.StatusCode = statusCode
	return builder
}

func (builder *ErrorBuilder) SetMessage(msg string) *ErrorBuilder {
	builder.Message = msg
	return builder
}

func (builder *ErrorBuilder) SetFieldError(field string, code string, details string) *ErrorBuilder {
	builder.Errors = append(builder.Errors, FieldError{field, code, details})
	return builder
}

func (builder *ErrorBuilder) SetMetaData(meta interface{}) *ErrorBuilder {
	builder.MetaData = meta
	return builder
}
func (builder *ErrorBuilder) SetError(err error) *ErrorBuilder {
	builder.OriginalError = err
	return builder
}

func (builder *ErrorBuilder) SetRequestId(requestId string) *ErrorBuilder {
	builder.RequestId = requestId
	return builder
}

func (builder *ErrorBuilder) Build() *ErrorResponse {
	return &ErrorResponse{
		RequestId: builder.RequestId,
		Message:   builder.Message,
		Errors:    builder.Errors,
	}
}

func NewInternalServerError(err error) *ErrorBuilder {
	return &ErrorBuilder{
		StatusCode:    http.StatusInternalServerError,
		Message:       http.StatusText(http.StatusInternalServerError),
		Errors:        nil,
		OriginalError: err,
		MetaData:      nil,
	}
}

func NewUnauthorized() *ErrorBuilder {
	return &ErrorBuilder{
		StatusCode: http.StatusUnauthorized,
		Message:    http.StatusText(http.StatusUnauthorized),
	}
}

func NewBadRequestError() *ErrorBuilder {
	return &ErrorBuilder{
		StatusCode:    http.StatusBadRequest,
		Message:       http.StatusText(http.StatusBadRequest),
		Errors:        nil,
		OriginalError: nil,
		MetaData:      nil,
	}
}

func NewNotFoundError() *ErrorBuilder {
	return &ErrorBuilder{
		StatusCode:    http.StatusNotFound,
		Message:       http.StatusText(http.StatusNotFound),
		Errors:        nil,
		OriginalError: nil,
		MetaData:      nil,
	}
}

func NewValidationError(err error) *ErrorBuilder {
	var fieldErrors []FieldError
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, validationError := range validationErrors {
			fieldErrors = append(fieldErrors, FieldError{
				Field:  validationError.Field(),
				Code:   validationError.Tag(),
				Detail: validationError.Error(),
			})
		}
	}
	return &ErrorBuilder{
		StatusCode:    http.StatusBadRequest,
		Message:       http.StatusText(http.StatusBadRequest),
		Errors:        fieldErrors,
		OriginalError: err,
		MetaData:      nil,
	}
}

func FromDatabase(err error) *ErrorBuilder {
	return &ErrorBuilder{
		StatusCode:    http.StatusBadRequest,
		Message:       http.StatusText(http.StatusBadRequest),
		OriginalError: err,
		MetaData:      nil,
	}
}
