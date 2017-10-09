package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd *cobra.Command = &cobra.Command{
	Use:   "jwt",
	Short: "A self-sufficient runtime for json-web-token auth",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	SetupRootCommand(RootCmd)
}
