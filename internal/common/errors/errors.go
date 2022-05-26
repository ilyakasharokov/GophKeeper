// customerrors пакет для работы с ошибками
package customerrors

import (
	"fmt"
)

// NewError функция получения новой ошибки
func NewError(err error, errorText string) error {
	return &CustomError{
		Err:       err,
		ErrorText: errorText,
	}
}

// ParseError функция по распаковки ошибки в строку
func ParseError(err error) string {
	switch e := err.(type) {
	case *CustomError:
		return e.ErrorText
	default:
		return "internal server error"
	}
}

// CustomError структура для хранения ошибок
type CustomError struct {
	Err       error
	ErrorText string
}

// Error возвращает текст ошибки
func (err *CustomError) Error() string {
	return fmt.Sprintf("%v", err.Err)
}

// Unwrap функция для распоковки ошибки
func (err *CustomError) Unwrap() error {
	return err.Err
}
