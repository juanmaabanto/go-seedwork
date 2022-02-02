package middleware

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/juanmaabanto/go-seedwork/seedwork/errors"
	"github.com/juanmaabanto/go-seedwork/seedwork/responses"
	"github.com/labstack/echo/v4"
)

type LoggerConfig struct {
	LoggerErrorFunc func(message string, trace string, username string, userAgent string) string
}

var DefaultLoggerConfig = LoggerConfig{
	LoggerErrorFunc: func(message string, trace string, username string, userAgent string) string {
		fmt.Println(message)

		return "Id autogenerado por defaultLoggerConfig"
	},
}

func Logger() echo.MiddlewareFunc {
	return LoggerWithConfig(DefaultLoggerConfig)
}

func LoggerWithConfig(config LoggerConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}

					var errorId string
					statusCode := getStatusCode(err)

					if statusCode == 500 {
						stack := make([]byte, 1<<10)
						length := runtime.Stack(stack, true)

						errorId = config.LoggerErrorFunc(err.Error(), string(stack[:length]), "test", c.Request().UserAgent())
					}

					response := responses.ErrorResponse{
						ErrorId: errorId,
						Message: getMessage(err, *c.Request()),
						Status:  statusCode,
						Title:   getTitle(err),
						Errors:  getErrors(err),
					}

					c.JSON(statusCode, response)
				}
			}()

			return next(c)
		}
	}
}

func getCustomMessage(request http.Request) string {
	switch request.Method {
	case http.MethodDelete:
		return "Se produjo un error al eliminar el recurso"
	case http.MethodGet:
		return "Se produjo un error obteniendo el recurso."
	case http.MethodPatch:
		return "Se produjo un error al intentar actualizar el recurso."
	case http.MethodPost:
		return "Se produjo un error al intentar crear el recurso."
	case http.MethodPut:
		return "Se produjo un error al intentar actualizar el recurso."
	default:
		return "Se produjo un error al consumir el servicio."
	}
}

func getErrors(err error) map[string]string {
	customErr, ok := err.(errors.ApplicationError)

	if !ok {
		return nil
	}

	if customErr.ErrorType() == errors.ErrorTypeValidation {
		return customErr.Errors()
	} else {
		return nil
	}
}

func getMessage(err error, request http.Request) string {
	customErr, ok := err.(errors.ApplicationError)

	if !ok {
		return getCustomMessage(request)
	}

	switch customErr.ErrorType() {
	case errors.ErrorTypeBadRequest, errors.ErrorTypeNotFound, errors.ErrorTypeValidation:
		return customErr.Error()
	default:
		return getCustomMessage(request)
	}
}

func getStatusCode(err error) int {
	customErr, ok := err.(errors.ApplicationError)

	if !ok {
		return http.StatusInternalServerError
	}

	switch customErr.ErrorType() {
	case errors.ErrorTypeBadRequest:
		return http.StatusBadRequest
	case errors.ErrorTypeNotFound:
		return http.StatusNotFound
	case errors.ErrorTypeValidation:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}

func getTitle(err error) string {
	customErr, ok := err.(errors.ApplicationError)

	if !ok {
		return "Server Error"
	} else {
		return customErr.Title()
	}
}
