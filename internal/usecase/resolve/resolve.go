package resolve

import (
	"fmt"
	"io"
	"os"

	"resodns/internal/app/ctx"
	"resodns/internal/pkg/console"
	"resodns/pkg/fileoperation"
	"resodns/pkg/shellexecutor"
)

type Service struct {
	Context *ctx.Ctx
	Options *ctx.ResolveOptions

	RequirementChecker RequirementChecker
	ResolverLoader     ResolverLoader
	WorkfileCreator    WorkfileCreator
	MassResolver       MassResolver
	ResultSaver        ResultSaver
	WildcardFilter     WildcardFilter

	workfiles   *Workfiles
	domainCount int
}

type RequirementChecker interface {
	Check(opt *ctx.ResolveOptions) error
}
type ResolverLoader interface {
	Load(ctx *ctx.Ctx, filename string) error
}
type WorkfileCreator interface {
	Create() (*Workfiles, error)
}
type MassResolver interface {
	Resolve(reader io.Reader, output string, total int, resolversFilename string, qps int) error
}
type ResultSaver interface {
	Save(workfiles *Workfiles, opt *ctx.ResolveOptions) error
}
type WildcardFilter interface {
	Filter(opt WildcardFilterOptions, totalCount int) (found int, roots []string, err error)
}

func NewService(c *ctx.Ctx, opt *ctx.ResolveOptions) *Service {
	return &Service{
		Context:            c,
		Options:            opt,
		RequirementChecker: NewDefaultRequirementChecker(shellexecutor.NewShellExecutor()),
		ResolverLoader:     NewDefaultResolverFileLoader(),
		WorkfileCreator:    NewDefaultWorkfileCreator(),
		MassResolver:       NewDefaultMassResolver(opt.BinPath),
		ResultSaver:        NewResultFileSaver(),
		WildcardFilter:     NewDefaultWildcardFilter(),
	}
}

func (s *Service) Initialize() error {
	if err := s.RequirementChecker.Check(s.Options); err != nil {
		return err
	}
	var err error
	if s.workfiles, err = s.WorkfileCreator.Create(); err != nil {
		return err
	}
	if err := s.prepareResolvers(); err != nil {
		return err
	}
	return nil
}

func (s *Service) Resolve() error {
	domainReader, err := s.createDomainReader()
	if err != nil {
		return err
	}
	if err = s.resolvePublic(domainReader); err != nil {
		return err
	}
	if err = s.filterWildcards(); err != nil {
		return err
	}
	if err = s.resolveTrusted(); err != nil {
		return err
	}
	if err = s.writeResults(); err != nil {
		return err
	}
	return nil
}

func (s *Service) Close(debug bool) {
	if debug {
		console.Printf("\nDebug files kept in: %s\n", s.workfiles.TempDirectory)
	} else if s.workfiles != nil {
		s.workfiles.Close()
	}
}

func (s *Service) prepareResolvers() error {
	if !s.Options.TrustedOnly {
		if err := fileoperation.Copy(s.Options.ResolverFile, s.workfiles.PublicResolvers); err != nil {
			return fmt.Errorf("unable to load public resolvers: %w", err)
		}
	}
	if err := s.ResolverLoader.Load(s.Context, s.Options.ResolverTrustedFile); err != nil {
		return fmt.Errorf("unable to load trusted resolvers: %w", err)
	}
	if err := fileoperation.WriteLines(s.Context.Options.TrustedResolvers, s.workfiles.TrustedResolvers); err != nil {
		return fmt.Errorf("unable to write trusted resolvers: %w", err)
	}
	return nil
}

func (s *Service) createDomainReader() (*DomainReader, error) {
	sourceReader, err := s.createDomainReaderSource()
	if err != nil {
		return nil, err
	}
	var domains []string
	if s.Options.Mode == ctx.Bruteforce {
		if domains, err = s.createDomainReaderDomainList(); err != nil {
			return nil, err
		}
	}
	var sanitizer DomainSanitizer
	if !s.Options.SkipSanitize {
		sanitizer = DefaultSanitizer
	}
	return NewDomainReader(sourceReader, domains, sanitizer), nil
}

func (s *Service) createDomainReaderSource() (io.ReadCloser, error) {
	if s.Context.Stdin != nil {
		return s.Context.Stdin, nil
	}
	var filename string
	if s.Options.Mode == ctx.Resolve {
		filename = s.Options.DomainFile
	} else {
		filename = s.Options.Wordlist
	}
	count, err := fileoperation.CountLines(filename)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	s.domainCount = count
	return file, nil
}

func (s *Service) createDomainReaderDomainList() ([]string, error) {
	var domains []string
	if s.Options.DomainFile != "" {
		var err error
		domains, err = fileoperation.ReadLines(s.Options.DomainFile)
		if err != nil {
			return nil, err
		}
	} else {
		domains = []string{s.Options.Domain}
	}
	s.domainCount = s.domainCount * len(domains)
	return domains, nil
}

func (s *Service) resolvePublic(reader *DomainReader) error {
	resolvers := s.workfiles.PublicResolvers
	ratelimit := s.Options.RateLimit
	resolverString := "public"
	if s.Options.TrustedOnly {
		resolvers = s.workfiles.TrustedResolvers
		ratelimit = s.Options.RateLimitTrusted
		resolverString = "trusted"
	}
	console.Printf("%sResolving domains with %s resolvers%s\n", console.ColorBrightWhite, resolverString, console.ColorReset)
	err := s.MassResolver.Resolve(reader, s.workfiles.MassdnsPublic, s.domainCount, resolvers, ratelimit)
	if err != nil {
		return fmt.Errorf("error resolving domains: %w", err)
	}
	console.Printf("\n")
	return nil
}

func (s *Service) filterWildcards() error {
	if s.Options.SkipWildcard {
		return s.parseCache(s.workfiles.MassdnsPublic, s.workfiles.Domains)
	}
	if err := s.parseCache(s.workfiles.MassdnsPublic, ""); err != nil {
		return err
	}
	console.Printf("%sDetecting wildcard root subdomains%s\n", console.ColorBrightWhite, console.ColorReset)
	opt := WildcardFilterOptions{
		CacheFilename:        s.workfiles.MassdnsPublic,
		DomainOutputFilename: s.workfiles.Domains,
		RootOutputFilename:   s.workfiles.WildcardRoots,
		Resolvers:            s.Context.Options.TrustedResolvers,
		QueriesPerSecond:     s.Options.RateLimitTrusted,
		ThreadCount:          s.Options.WildcardThreads,
		ResolveTestCount:     s.Options.WildcardTests,
		BatchSize:            s.Options.WildcardBatchSize,
	}
	found, roots, err := s.WildcardFilter.Filter(opt, s.domainCount)
	if err != nil {
		return fmt.Errorf("unable to filter wildcard domains: %w", err)
	}
	if len(roots) > 0 {
		console.Printf("\n%sFound %s%d%s wildcard roots:%s\n", console.ColorBrightWhite, console.ColorBrightGreen, len(roots), console.ColorBrightWhite, console.ColorReset)
		for _, root := range roots {
			console.Printf("*.%s\n", root)
		}
	}
	s.domainCount = found
	console.Printf("\n")
	return nil
}

func (s *Service) parseCache(cacheFilename string, domainFilename string) error {
	var domainFile *os.File
	cacheFile, err := os.Open(cacheFilename)
	if err != nil {
		return err
	}
	if domainFilename != "" {
		if domainFile, err = os.Create(domainFilename); err != nil {
			return err
		}
	}
	cacheReader := NewCacheReader(cacheFile)
	if s.domainCount, err = cacheReader.Read(domainFile, nil, 0); err != nil {
		return err
	}
	if domainFilename != "" {
		domainFile.Sync()
		cacheReader.Close()
		return domainFile.Close()
	}
	return cacheReader.Close()
}

func (s *Service) resolveTrusted() error {
	if s.Options.SkipValidation {
		return nil
	}
	domainFile, err := os.Open(s.workfiles.Domains)
	if err != nil {
		return nil
	}
	defer domainFile.Close()
	console.Printf("%sValidating domains against trusted resolvers%s\n", console.ColorBrightWhite, console.ColorReset)
	err = s.MassResolver.Resolve(domainFile, s.workfiles.MassdnsTrusted, s.domainCount, s.workfiles.TrustedResolvers, s.Options.RateLimitTrusted)
	if err != nil {
		return fmt.Errorf("error resolving domains: %w", err)
	}
	console.Printf("\n")
	return s.parseCache(s.workfiles.MassdnsTrusted, s.workfiles.Domains)
}

func (s *Service) writeResults() error {
	if s.domainCount > 0 {
		console.Printf("%sFound %s%d%s valid domains:%s\n", console.ColorBrightWhite, console.ColorBrightGreen, s.domainCount, console.ColorBrightWhite, console.ColorReset)
	} else {
		console.Printf("\n%sNo valid domains remaining.%s\n", console.ColorBrightWhite, console.ColorReset)
	}
	if err := fileoperation.Cat([]string{s.workfiles.Domains}, os.Stdout); err != nil {
		return fmt.Errorf("unable to read domain file: %w", err)
	}
	if err := s.ResultSaver.Save(s.workfiles, s.Options); err != nil {
		return fmt.Errorf("unable to save results: %w", err)
	}
	return nil
}
