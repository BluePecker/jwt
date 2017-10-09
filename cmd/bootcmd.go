package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd *cobra.Command = &cobra.Command{
	Use:   "jwt",
	Short: "Jwt auth server",
	Long:  "User login information verification service",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	SilenceErrors: true,
}

func init() {
	SetupRootCommand(RootCmd)
}
