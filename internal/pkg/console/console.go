package console

import (
	"fmt"
	"io"
	"os"
)

const (
	colorMeta        = ColorBrightWhite
	colorMessage     = ColorCyan
	colorMessageText = ColorWhite
	colorSuccess     = ColorBrightGreen
	colorSuccessText = ColorWhite
	colorWarning     = ColorBrightYellow
	colorWarningText = ColorWhite
	colorError       = ColorRed
	colorErrorText   = ColorRed
)

var (
	Output      io.Writer = os.Stderr
	ExitHandler            = os.Exit
)

func Message(format string, a ...interface{}) {
	msg := fmt.Sprintf("%s[%s*%s]%s %s%s\n", colorMeta, colorMessage, colorMeta, colorMessageText, format, ColorReset)
	fmt.Fprintf(Output, msg, a...)
}

func Success(format string, a ...interface{}) {
	msg := fmt.Sprintf("%s[%s+%s]%s %s%s\n", colorMeta, colorSuccess, colorMeta, colorSuccessText, format, ColorReset)
	fmt.Fprintf(Output, msg, a...)
}

func Warning(format string, a ...interface{}) {
	msg := fmt.Sprintf("%s[%s!%s]%s %s%s\n", colorMeta, colorWarning, colorMeta, colorWarningText, format, ColorReset)
	fmt.Fprintf(Output, msg, a...)
}

func Error(format string, a ...interface{}) {
	msg := fmt.Sprintf("%s[%sX%s]%s %s%s\n", colorMeta, colorError, colorMeta, colorErrorText, format, ColorReset)
	fmt.Fprintf(Output, msg, a...)
}

func Printf(format string, a ...interface{}) {
	fmt.Fprintf(Output, format, a...)
}

func Fatal(format string, a ...interface{}) {
	fmt.Fprintf(Output, ColorRed+format+ColorReset, a...)
	ExitHandler(-1)
}
