package programbanner

import (
	"fmt"
	"strings"

	"resodns/internal/app/ctx"
	"resodns/internal/pkg/console"
)

type Service struct {
	ctx *ctx.Ctx
}

func NewService(c *ctx.Ctx) Service {
	return Service{ctx: c}
}

func (s Service) Print() {
	version := s.ctx.ProgramVersion
	if s.ctx.GitBranch != "" {
		version = fmt.Sprintf("%s-%s", s.ctx.GitBranch, s.ctx.GitRevision)
	}
	padding := strings.Repeat(" ", 34-len(version)-len(s.ctx.ProgramName))
	console.Printf(console.ColorBrightBlue)
	console.Printf("%s%s%s %s%s\n", padding, console.ColorBrightCyan, s.ctx.ProgramName, console.ColorBrightBlue, version)
	console.Printf("\n%sMass DNS resolution and subdomain bruteforce with wildcard filtering\n%s\n", console.ColorBrightWhite, console.ColorReset)
}

func (s Service) PrintWithResolveOptions(opts *ctx.ResolveOptions) {
	s.Print()
	console.Printf(console.ColorBrightWhite + "------------------------------------------------------------\n" + console.ColorReset)
	def := ctx.DefaultResolveOptions()
	var file string
	if s.ctx.Stdin != nil {
		file = "stdin"
	} else {
		if opts.Mode == ctx.Bruteforce {
			file = opts.Wordlist
		} else {
			file = opts.DomainFile
		}
	}
	cL, cSkip, cVal, cTick, cTickW := console.ColorBrightWhite, console.ColorBrightYellow, console.ColorWhite, console.ColorBrightBlue, console.ColorBrightGreen
	tick := fmt.Sprintf("%s[%s+%s]", cL, cTick, cL)
	tickW := fmt.Sprintf("%s[%s+%s]", cL, cTickW, cL)
	if opts.Mode == ctx.Bruteforce {
		console.Printf("%s Mode                 :%s bruteforce\n", tick, cVal)
		if opts.DomainFile != "" {
			console.Printf("%s Domains              :%s %s\n", tick, cVal, opts.DomainFile)
		} else {
			console.Printf("%s Domain               :%s %s\n", tick, cVal, opts.Domain)
		}
		console.Printf("%s Wordlist             :%s %s\n", tick, cVal, file)
	} else {
		console.Printf("%s Mode                 :%s resolve\n", tick, cVal)
		console.Printf("%s File                 :%s %s\n", tick, cVal, file)
	}
	if opts.TrustedOnly {
		console.Printf("%s Trusted Only         :%s true\n", tick, cVal)
	}
	if !opts.TrustedOnly {
		console.Printf("%s Resolvers            :%s %s\n", tick, cVal, opts.ResolverFile)
	}
	if opts.ResolverTrustedFile != "" {
		console.Printf("%s Trusted Resolvers    :%s %s\n", tick, cVal, opts.ResolverTrustedFile)
	}
	if !opts.TrustedOnly {
		rate := "unlimited"
		if opts.RateLimit != 0 {
			rate = fmt.Sprintf("%d qps", opts.RateLimit)
		}
		console.Printf("%s Rate Limit           :%s %s\n", tick, cVal, rate)
	}
	console.Printf("%s Rate Limit (Trusted) :%s %d qps\n", tick, cVal, opts.RateLimitTrusted)
	console.Printf("%s Wildcard Threads     :%s %d\n", tick, cVal, opts.WildcardThreads)
	console.Printf("%s Wildcard Tests       :%s %d\n", tick, cVal, opts.WildcardTests)
	if opts.WildcardBatchSize != def.WildcardBatchSize {
		console.Printf("%s Wildcard Batch Size  :%s %d\n", tick, cVal, opts.WildcardBatchSize)
	}
	if opts.WriteDomainsFile != "" {
		console.Printf("%s Write Domains        :%s %s\n", tickW, cVal, opts.WriteDomainsFile)
	}
	if opts.WriteMassdnsFile != "" {
		console.Printf("%s Write Massdns        :%s %s\n", tickW, cVal, opts.WriteMassdnsFile)
	}
	if opts.WriteWildcardsFile != "" {
		console.Printf("%s Write Wildcards      :%s %s\n", tickW, cVal, opts.WriteWildcardsFile)
	}
	if opts.SkipSanitize {
		console.Printf("%s[+] Skip Sanitize\n", cSkip)
	}
	if opts.SkipWildcard {
		console.Printf("%s[+] Skip Wildcard Detection\n", cSkip)
	}
	if !opts.TrustedOnly && opts.SkipValidation {
		console.Printf("%s[+] Skip Validation\n", cSkip)
	}
	console.Printf(console.ColorBrightWhite + "------------------------------------------------------------\n" + console.ColorReset)
	console.Printf("\n")
}
