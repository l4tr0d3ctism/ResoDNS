package resolve

import (
	"io"

	"resodns/internal/pkg/console"
	"resodns/pkg/massdns"
	"resodns/pkg/progressbar"
)

type DefaultMassResolver struct {
	massdns *massdns.Resolver
}

func NewDefaultMassResolver(binPath string) *DefaultMassResolver {
	return &DefaultMassResolver{massdns: massdns.NewResolver(binPath)}
}

func (m *DefaultMassResolver) Resolve(r io.Reader, output string, total int, resolversFilename string, qps int) error {
	var template string
	if total == 0 {
		template = "Processed: {{ current }} Rate: {{ rate }} Elapsed: {{ time }}"
	} else {
		template = "[ETA {{ eta }}] {{ bar }} {{ current }}/{{ total }} rate: {{ rate }} qps (time: {{ time }})"
	}
	bar := progressbar.New(m.updateProgressBar, int64(total), progressbar.WithTemplate(template), progressbar.WithWriter(console.Output))
	bar.Start()
	err := m.massdns.Resolve(r, output, resolversFilename, qps)
	bar.Stop()
	return err
}

func (m *DefaultMassResolver) updateProgressBar(bar *progressbar.ProgressBar) {
	bar.SetCurrent(int64(m.massdns.Current()))
}
