package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd *cobra.Command = &cobra.Command{

}

func init() {
	SetupRootCommand(RootCmd)
}
