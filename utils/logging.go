package utils

import (
	"fmt"
	"golang.org/x/exp/slog"
)

func Infof(format string, args ...any) {
	slog.Default().Info(fmt.Sprintf(format, args...))
}
