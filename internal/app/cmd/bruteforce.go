package cmd

import (
	"os"

	"resodns/internal/app"
	"resodns/internal/app/ctx"
	"resodns/internal/usecase/programbanner"
	"resodns/internal/usecase/resolve"
	"github.com/spf13/cobra"
)

func newCmdBruteforce() *cobra.Command {
	cmdBruteforce := &cobra.Command{
		Use:   "bruteforce <wordlist> domain [flags]\n  resodns bruteforce <wordlist> -d domains.txt [flags]",
		Short: "Bruteforce subdomains using a wordlist",
		Long: `Bruteforce subdomains with a wordlist. The <wordlist> argument can be omitted
if the wordlist is read from stdin (e.g. cat wordlist.txt | resodns bruteforce - example.com).`,
		RunE: runBruteforce,
	}
	cmdBruteforce.Flags().StringVarP(&resolveOptions.DomainFile, "domains", "d", resolveOptions.DomainFile, "file containing domains to bruteforce")
	cmdBruteforce.Flags().AddFlagSet(resolveFlags)
	cmdBruteforce.Flags().SortFlags = false
	return cmdBruteforce
}

func runBruteforce(cmd *cobra.Command, args []string) error {
	parseBruteforceArgs(args)
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

func parseBruteforceArgs(args []string) {
	if app.HasStdin() {
		context.Stdin = os.Stdin
		if len(args) >= 1 && resolveOptions.DomainFile == "" {
			resolveOptions.Domain = args[0]
		}
	} else {
		if len(args) == 1 {
			resolveOptions.Wordlist = args[0]
		} else if len(args) >= 2 {
			resolveOptions.Wordlist = args[0]
			resolveOptions.Domain = args[1]
		}
	}
	resolveOptions.Mode = ctx.Bruteforce
}
