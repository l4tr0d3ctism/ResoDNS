package main

import (
	"fmt"
	"os"

	"github.com/l4tr0d3ctism/ResoDNS/internal/app/cmd"
	"github.com/l4tr0d3ctism/ResoDNS/internal/app/ctx"
	"github.com/l4tr0d3ctism/ResoDNS/internal/app/errmsg"
)

var exitHandler func(int) = os.Exit

func main() {
	ctx := ctx.NewCtx()

	if err := cmd.Execute(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s%s\n", errmsg.Prefix, err)
		exitHandler(1)
	}
}
