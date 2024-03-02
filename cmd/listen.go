/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen some dir changing and exec cmd on it",
	Long: `
		fsd listen data quicktype --src data --src-lang json --lang dart --out flutter/types/lib/src/types_base.dart --coders-in-class --from-map --final-props --required-props --copy-with
		fsd listen data quicktype --src data --src-lang json --lang go --out go/src/types/types.go --package types
	`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.SetCommandPath(cmd.CommandPath())
		config.Bind("listenDir")
		config.Bind("execCmd")
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmds := []string{
			fmt.Sprintf("quicktype --src %s --src-lang json --lang dart --out %s --coders-in-class --from-map --final-props --required-props --copy-with", config.GetString("project_root")+config.GetString("listenDirToGenerate"), config.GetString("project_root")+config.GetString("generatedDartDir")),
			fmt.Sprintf("quicktype --src %s --src-lang json --lang go --out %s --package types", config.GetString("project_root")+config.GetString("listenDirToGenerate"), config.GetString("project_root")+config.GetString("generatedGotDir")),
		}
		execOnChange(config.ConfigPath+"data", cmds)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		config.WriteCfgToFile()
	},
}

func init() {
	MainCmd.AddCommand(listenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	listenCmd.PersistentFlags().StringP("listenDir", "l", "./input", "dir for listening changes")
	listenCmd.PersistentFlags().StringP("execCmd", "e", "./output", "cmd to execute on change")
}
