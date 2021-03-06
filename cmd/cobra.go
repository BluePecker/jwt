package cmd

import (
	"strings"
	"fmt"
	"github.com/spf13/cobra"
)

func usageTemplate() string {
	return `Usage:	{{if not .HasSubCommands}}{{.Use}} [OPTIONS] COMMAND [arg...]{{end}}{{if .HasSubCommands}}{{ .CommandPath}} COMMAND{{end}}

{{ .Short | trim }}{{if gt .Aliases 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{ .Example }}{{end}}{{if .HasFlags}}

Options:
{{.Flags.FlagUsages | trimRightSpace}}{{end}}{{ if .HasAvailableSubCommands}}

Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasSubCommands }}

Run '{{.CommandPath}} COMMAND --help' for more information on a command.{{end}}
`
}

func helpTemplate() string {
	return "\n{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}"
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
	cmd.SetHelpCommand(&cobra.Command{
		Use:   "help [command]",
		Short: "Help about the command",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// todo
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			// todo
		},
		RunE: func(c *cobra.Command, args []string) error {
			cmd, args, err := c.Root().Find(args)
			if cmd == nil || err != nil || len(args) > 0 {
				return fmt.Errorf("unknown help topic: %v", strings.Join(args, " "))
			}
			cmd.HelpFunc()(cmd, args)
			return nil
		},
	})
	cmd.SetFlagErrorFunc(FlagErrorFunc)

	cmd.PersistentFlags().BoolP("help", "h", false, "print usage")
	cmd.PersistentFlags().MarkShorthandDeprecated("help", "please use --help")
}
