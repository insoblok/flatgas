package cmd

import "github.com/spf13/cobra"

func GetNodeCommand() *cobra.Command {
	nodeCmd := &cobra.Command{
		Use:   "node",
		Short: "Node-related commands",
	}

	nodeCmd.AddCommand(GetNodeInfoCommand())
	return nodeCmd
}
