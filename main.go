package main

import (
	"fmt"
	"os"

	"resodns/internal/app/cmd"
	"resodns/internal/app/ctx"
	"resodns/internal/app/errmsg"
)

var exitHandler func(int) = os.Exit

func main() {
	ctx := ctx.NewCtx()

	if err := cmd.Execute(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s%s\n", errmsg.Prefix, err)
		exitHandler(1)
	}
}
