package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd *cobra.Command = &cobra.Command{
	Use:   "jwt",
	Long:  "User login information verification service",
	Short: "Jwt auth server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	SilenceErrors: true,
}

func init() {
	SetupRootCommand(RootCmd)
}
