package errs

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type HTTPError struct {
	Code    int `json:"code"`
	Message any `json:"message"`
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("http error: %d %v", e.Code, e.Message)
}

func NewHTTPError(code int, err error) HTTPError {
	return HTTPError{
		Code:    code,
		Message: err.Error(),
	}
}

// BadRequest accepts optional custom message, defaults to "bad request"
func BadRequest(msg ...string) HTTPError {
	message := "bad request"
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	return NewHTTPError(http.StatusBadRequest, fmt.Errorf("%s", message))
}

// Unauthorized accepts optional custom message, defaults to "unauthorized"
func Unauthorized(msg ...string) HTTPError {
	message := "unauthorized"
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	return NewHTTPError(http.StatusUnauthorized, errors.New(message))
}

// NotFound with flexible parameters
func NotFound(msg ...string) HTTPError {
	if len(msg) == 0 {
		return NewHTTPError(http.StatusNotFound, errors.New("resource not found"))
	}

	// If only one argument, use it as the message
	if len(msg) == 1 {
		return NewHTTPError(http.StatusNotFound, errors.New(msg[0]))
	}

	// If three arguments, format as: "title with key='value' not found"
	if len(msg) == 3 {
		return NewHTTPError(http.StatusNotFound,
			fmt.Errorf("%s with %s='%s' not found", msg[0], msg[1], msg[2]))
	}

	// Otherwise join all messages
	return NewHTTPError(http.StatusNotFound, fmt.Errorf("%v", msg))
}

// Conflict with flexible parameters
func Conflict(msg ...string) HTTPError {
	if len(msg) == 0 {
		return NewHTTPError(http.StatusConflict, fmt.Errorf("resource conflict"))
	}

	// If only one argument, use it as the message
	if len(msg) == 1 {
		return NewHTTPError(http.StatusConflict, fmt.Errorf(msg[0]))
	}

	// If three arguments, format as: "title with key='value' already exists"
	if len(msg) == 3 {
		return NewHTTPError(http.StatusConflict,
			fmt.Errorf("%s with %s='%s' already exists", msg[0], msg[1], msg[2]))
	}

	// Otherwise join all messages
	return NewHTTPError(http.StatusConflict, fmt.Errorf("%v", msg))
}

// InvalidRequestData for validation errors
func InvalidRequestData(errors map[string]string) HTTPError {
	return HTTPError{
		Code:    http.StatusBadRequest,
		Message: errors,
	}
}

// InvalidJSON accepts optional custom message
func InvalidJSON(msg ...string) HTTPError {
	message := "invalid json"
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	return NewHTTPError(http.StatusBadRequest, errors.New(message))
}

// InternalServerError accepts optional custom message
func InternalServerError(msg ...string) HTTPError {
	message := "internal server error"
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	return NewHTTPError(http.StatusInternalServerError, errors.New(message))
}

// Forbidden accepts optional custom message
func Forbidden(msg ...string) HTTPError {
	message := "forbidden"
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	return NewHTTPError(http.StatusForbidden, errors.New(message))
}

// UnprocessableEntity for single message (not validation errors)
func UnprocessableEntity(msg ...string) HTTPError {
	message := "unprocessable entity"
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	return NewHTTPError(http.StatusUnprocessableEntity, errors.New(message))
}

// ErrorHandler remains the same
func ErrorHandler(c *fiber.Ctx, err error) error {
	var httpErr HTTPError
	if castedErr, ok := err.(HTTPError); ok {
		httpErr = castedErr
	} else {
		httpErr = InternalServerError()
	}

	slog.Error("HTTP API error", "err", err.Error(), "method", c.Method(), "path", c.Path())
	return c.Status(httpErr.Code).JSON(httpErr)
}
