package logx

import (
	"fmt"

	logger "github.com/yourpov/logrite"
)

// Err returns the error using logrite
func Err(format string, args ...any) error {
	err := fmt.Errorf(format, args...)
	logger.Error("%v", err)
	return err
}
