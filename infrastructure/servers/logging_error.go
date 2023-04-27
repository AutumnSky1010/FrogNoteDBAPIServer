package servers

import "fmt"

// LoggingError は、ロギング時のエラーを表現する構造体です。
type LoggingError struct {
	errorMessage  string
	internalError string
}

func (l LoggingError) Error() string {
	err := fmt.Sprintf("LoggingError: %s", l.errorMessage)
	if l.internalError != "" {
		err += fmt.Sprintf("\n\tinternalError: %s", l.internalError)
	}
	return err
}

func NewLoggingError(errorMessage string, internalError string) LoggingError {
	return LoggingError{errorMessage: errorMessage, internalError: internalError}
}
