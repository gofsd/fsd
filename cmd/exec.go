/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/gofsd/fsd/pkg/util"
	"github.com/spf13/cobra"
)

func Exec(rootCmd *cobra.Command) *cobra.Command {

	var execCmd = &cobra.Command{
		Use:   "exec",
		Short: "Create tree",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			o, e := util.Exec(args[0], 2)
			fmt.Print(string(o))
			return e
		},
	}
	execCmd.Args = cobra.MatchAll(cobra.MinimumNArgs(1), cobra.MaximumNArgs(1))

	rootCmd.AddCommand(execCmd)

	return execCmd
}

func init() {
	Exec(MainCmd)
}
