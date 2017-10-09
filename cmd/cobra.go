package cmd

import (
	"strings"
	"fmt"
	"github.com/spf13/cobra"
)

func usageTemplate() string {
	return `Usage:{{if .Runnable}}{{if .HasAvailableFlags}}
  {{appendIfNotPresent .UseLine "[OPTIONS] COMMAND [arg...]"}}{{else}}{{.UseLine}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
  {{ .CommandPath}} [command]
  {{end}}{{if gt .Aliases 0}}
Aliases:{{.NameAndAliases}}
{{end}}{{if .HasExample}}
Examples:{{ .Example }}
{{end}}{{ if .HasAvailableLocalFlags}}
Options:
{{.LocalFlags.FlagUsages | trimRightSpace}}
{{end}}{{ if .HasAvailableSubCommands}}
Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}
{{end}}{{ if .HasAvailableInheritedFlags}}
Global Flags:
{{.InheritedFlags.FlagUsages | trimRightSpace}}{{end}}{{if .HasHelpSubCommands}}
Additional help topics:{{range .Commands}}{{if .IsHelpCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableSubCommands }}
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}

func helpCommand() *cobra.Command {
	return &cobra.Command{
		Use:               "help [command]",
		Short:             "Help about the command",
		PersistentPreRun:  func(cmd *cobra.Command, args []string) {},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {},
		RunE: func(c *cobra.Command, args []string) error {
			cmd, args, e := c.Root().Find(args)
			if cmd == nil || e != nil || len(args) > 0 {
				return fmt.Errorf("unknown help topic: %v", strings.Join(args, " "))
			}

			helpFunc := cmd.HelpFunc()
			helpFunc(cmd, args)
			return nil
		},
	}
}

func helpTemplate() string {
	return `
{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
}

func FlagErrorFunc(cmd *cobra.Command, err error) error {
	if err == nil {
		return nil
	}

	usage := ""
	if cmd.HasSubCommands() {
		usage = "\n\n" + cmd.UsageString()
	}
	return fmt.Errorf("%s\nSee '%s --help'.%s", err, cmd.CommandPath(), usage)
}

func SetupRootCommand(cmd *cobra.Command) {
	cmd.SetUsageTemplate(usageTemplate())
	cmd.SetHelpTemplate(helpTemplate())
	cmd.SetHelpCommand(helpCommand())
	cmd.SetFlagErrorFunc(FlagErrorFunc)

	cmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
	cmd.PersistentFlags().MarkShorthandDeprecated("help", "please use --help")
}