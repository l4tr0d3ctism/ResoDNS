package resolve

import (
	"fmt"

	"resodns/internal/app/ctx"
)

type Executor interface {
	Shell(name string, arg ...string) error
}

type DefaultRequirementChecker struct {
	executor Executor
}

func NewDefaultRequirementChecker(executor Executor) DefaultRequirementChecker {
	return DefaultRequirementChecker{executor: executor}
}

func (c DefaultRequirementChecker) Check(opt *ctx.ResolveOptions) error {
	if err := c.executor.Shell(opt.BinPath, "--help"); err != nil {
		fmt.Printf("Unable to execute massdns. Make sure it is present and that the\n")
		fmt.Printf("path to the binary is added to the PATH environment variable.\n\n")
		fmt.Printf("Alternatively, specify the path to massdns using --bin\n\n")
		return fmt.Errorf("unable to execute massdns: %w", err)
	}
	return nil
}
