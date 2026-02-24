package ctx

import (
	"os"

	"resodns/internal/app"
)

type Ctx struct {
	ProgramName    string
	ProgramVersion string
	ProgramTagline string
	GitBranch      string
	GitRevision    string
	Options        *GlobalOptions
	Stdin          *os.File
}

func NewCtx() *Ctx {
	return &Ctx{
		ProgramName:    app.AppName,
		ProgramVersion: app.AppVersion,
		ProgramTagline: app.AppDesc,
		GitBranch:      app.GitBranch,
		GitRevision:    app.GitRevision,
		Options:        DefaultGlobalOptions(),
	}
}
