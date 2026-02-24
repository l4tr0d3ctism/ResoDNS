package ctx

import (
	"errors"
	"os/user"
	"path/filepath"

	"resodns/internal/app"
	"resodns/pkg/fileoperation"
)

type ResolveMode int

const (
	Resolve    ResolveMode = iota
	Bruteforce
)

var (
	ErrNoDomain   = errors.New("no domain specified")
	ErrNoWordlist = errors.New("no wordlist specified")
)

type GlobalOptions struct {
	TrustedResolvers []string
	Quiet            bool
	Debug            bool
}

func DefaultGlobalOptions() *GlobalOptions {
	return &GlobalOptions{
		TrustedResolvers: []string{"8.8.8.8", "8.8.4.4"},
		Quiet:            false,
		Debug:            false,
	}
}

type ResolveOptions struct {
	BinPath              string
	ResolverFile         string
	ResolverTrustedFile  string
	TrustedOnly          bool
	RateLimit            int
	RateLimitTrusted     int
	WildcardThreads      int
	WildcardTests        int
	WildcardBatchSize    int
	SkipSanitize         bool
	SkipWildcard         bool
	SkipValidation       bool
	WriteDomainsFile     string
	WriteMassdnsFile     string
	WriteWildcardsFile   string
	Mode                 ResolveMode
	Domain               string
	Wordlist             string
	DomainFile           string
}

func DefaultResolveOptions() *ResolveOptions {
	resolversPath := "resolvers.txt"
	trustedResolversPath := ""
	if !fileoperation.FileExists(resolversPath) {
		usr, err := user.Current()
		if err == nil {
			resolversPath = filepath.Join(usr.HomeDir, ".config", "resodns", "resolvers.txt")
			trustedResolversPath = filepath.Join(usr.HomeDir, ".config", "resodns", "resolvers-trusted.txt")
			if !fileoperation.FileExists(trustedResolversPath) {
				trustedResolversPath = ""
			}
		}
	}
	return &ResolveOptions{
		BinPath:             "massdns",
		ResolverFile:        resolversPath,
		ResolverTrustedFile: trustedResolversPath,
		TrustedOnly:         false,
		RateLimit:           0,
		RateLimitTrusted:    500,
		WildcardThreads:     100,
		WildcardTests:       3,
		WildcardBatchSize:    0,
		SkipSanitize:        false,
		SkipWildcard:        false,
		SkipValidation:      false,
		WriteDomainsFile:    "",
		WriteMassdnsFile:    "",
		WriteWildcardsFile:  "",
		Mode:                Resolve,
		Domain:              "",
		Wordlist:            "",
		DomainFile:          "",
	}
}

func (o *ResolveOptions) Validate() error {
	if o.TrustedOnly {
		o.SkipValidation = true
	}
	if o.Mode == Bruteforce {
		if o.Domain == "" && o.DomainFile == "" {
			return ErrNoDomain
		}
		if o.Wordlist == "" && !app.HasStdin() {
			return ErrNoWordlist
		}
	}
	return nil
}
