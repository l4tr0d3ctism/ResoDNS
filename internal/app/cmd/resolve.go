package cmd

import (
	"errors"
	"fmt"
	"os"

	"resodns/internal/app"
	"resodns/internal/app/ctx"
	"resodns/internal/usecase/programbanner"
	"resodns/internal/usecase/resolve"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	resolveFlags   *pflag.FlagSet
	resolveOptions *ctx.ResolveOptions
)

func newCmdResolve() *cobra.Command {
	resolveOptions = ctx.DefaultResolveOptions()
	cmdResolve := &cobra.Command{
		Use:   "resolve <file> [flags]",
		Short: "Resolve a list of domains",
		Long: `Resolve takes a file or stdin containing domains and performs DNS queries.
Use - or omit <file> to read from stdin.`,
		Args: cobra.MinimumNArgs(0),
		RunE: runResolve,
	}
	resolveFlags = pflag.NewFlagSet("resolve", pflag.ExitOnError)
	resolveFlags.StringVarP(&resolveOptions.BinPath, "bin", "b", resolveOptions.BinPath, "path to massdns binary")
	resolveFlags.IntVarP(&resolveOptions.RateLimit, "rate-limit", "l", resolveOptions.RateLimit, "queries per second for public resolvers (0 = unlimited)")
	resolveFlags.IntVar(&resolveOptions.RateLimitTrusted, "rate-limit-trusted", resolveOptions.RateLimitTrusted, "queries per second for trusted resolvers")
	resolveFlags.StringVarP(&resolveOptions.ResolverFile, "resolvers", "r", resolveOptions.ResolverFile, "file containing public resolvers")
	resolveFlags.StringVar(&resolveOptions.ResolverTrustedFile, "resolvers-trusted", resolveOptions.ResolverTrustedFile, "file containing trusted resolvers")
	resolveFlags.IntVarP(&resolveOptions.WildcardThreads, "threads", "t", resolveOptions.WildcardThreads, "threads for wildcard filtering")
	resolveFlags.IntVarP(&resolveOptions.WildcardTests, "wildcard-tests", "n", resolveOptions.WildcardTests, "tests for DNS load balancing detection")
	resolveFlags.IntVar(&resolveOptions.WildcardBatchSize, "wildcard-batch", resolveOptions.WildcardBatchSize, "subdomains per wildcard batch (0 = unlimited)")
	resolveFlags.StringVarP(&resolveOptions.WriteDomainsFile, "write", "w", resolveOptions.WriteDomainsFile, "write found domains to file")
	resolveFlags.StringVar(&resolveOptions.WriteMassdnsFile, "write-massdns", resolveOptions.WriteMassdnsFile, "write massdns output to file")
	resolveFlags.StringVar(&resolveOptions.WriteWildcardsFile, "write-wildcards", resolveOptions.WriteWildcardsFile, "write wildcard roots to file")
	resolveFlags.BoolVar(&resolveOptions.SkipSanitize, "skip-sanitize", resolveOptions.SkipSanitize, "do not sanitize domains")
	resolveFlags.BoolVar(&resolveOptions.SkipWildcard, "skip-wildcard-filter", resolveOptions.SkipWildcard, "skip wildcard detection")
	resolveFlags.BoolVar(&resolveOptions.SkipValidation, "skip-validation", resolveOptions.SkipValidation, "skip validation with trusted resolvers")

	cmdResolve.Flags().AddFlagSet(resolveFlags)
	cmdResolve.Flags().SortFlags = false
	return cmdResolve
}

func runResolve(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		if !app.HasStdin() {
			fmt.Fprint(os.Stderr, cmd.UsageString())
			return errors.New("requires a list of domains to resolve (file or stdin)")
		}
		context.Stdin = os.Stdin
	} else if args[0] == "-" {
		context.Stdin = os.Stdin
	} else {
		resolveOptions.DomainFile = args[0]
	}
	resolveOptions.Mode = ctx.Resolve
	if err := resolveOptions.Validate(); err != nil {
		return err
	}
	bannerService := programbanner.NewService(context)
	resolveService := resolve.NewService(context, resolveOptions)
	if err := resolveService.Initialize(); err != nil {
		return err
	}
	defer resolveService.Close(context.Options.Debug)
	bannerService.PrintWithResolveOptions(resolveOptions)
	return resolveService.Resolve()
}
