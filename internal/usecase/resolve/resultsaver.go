package resolve

import (
	"github.com/l4tr0d3ctism/ResoDNS/internal/app/ctx"
	"github.com/l4tr0d3ctism/ResoDNS/pkg/fileoperation"
)

type ResultFileSaver struct {
	fileCopy func(src, dest string) error
}

func NewResultFileSaver() *ResultFileSaver {
	return &ResultFileSaver{fileCopy: fileoperation.Copy}
}

func (s *ResultFileSaver) Save(workfiles *Workfiles, opt *ctx.ResolveOptions) error {
	if opt.WriteDomainsFile != "" {
		if err := s.fileCopy(workfiles.Domains, opt.WriteDomainsFile); err != nil {
			return err
		}
	}
	if opt.WriteMassdnsFile != "" {
		if err := s.fileCopy(workfiles.MassdnsPublic, opt.WriteMassdnsFile); err != nil {
			return err
		}
	}
	if opt.WriteWildcardsFile != "" {
		if err := s.fileCopy(workfiles.WildcardRoots, opt.WriteWildcardsFile); err != nil {
			return err
		}
	}
	return nil
}
