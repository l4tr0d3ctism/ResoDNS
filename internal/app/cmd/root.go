package cmd

import (
	"resodns/internal/app/ctx"
	"github.com/spf13/cobra"
)

var context *ctx.Ctx

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func newCmdRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   context.ProgramName,
		Short: context.ProgramTagline,
		Long:  context.ProgramName + " " + context.ProgramVersion + "\n\nMass DNS resolution and subdomain bruteforce with wildcard filtering.",
		Example: `  resodns resolve domains.txt
  resodns bruteforce wordlist.txt example.com --resolvers public.txt
  cat domains.txt | resodns resolve
  cat wordlist.txt | resodns bruteforce - example.com`,
		Version: context.ProgramVersion,
	}
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Flags().SortFlags = false

	rootCmd.AddCommand(newCmdResolve())
	rootCmd.AddCommand(newCmdBruteforce())
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.SilenceErrors = true
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) { cmd.SilenceUsage = true }
	return rootCmd
}

func Execute(c *ctx.Ctx) error {
	context = c
	return newCmdRoot().Execute()
}
