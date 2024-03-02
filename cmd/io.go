/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/gofsd/fsd/pkg/pipe"
	"github.com/spf13/cobra"
)

// ioCmd represents the io command
var ioCmd = &cobra.Command{
	Use:   "io",
	Short: "global IO command",
	Long:  `IO for everything`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.SetCommandPath(cmd.CommandPath())
		config.Bind("top-words-by-google")
		config.Bind("piplineName")
		config.Bind("show-top20k")

	},
	Run: func(cmd *cobra.Command, args []string) {
		switch config.GetString("piplineName") {
		case "moveTop20ToDB":
			pipe.MoveTop20kWordListToDB(config.GetString("top-words-by-google"), config.GetString("db_path"))
		case "test":
			pipe.Test(config.GetString("db_path"))
		}
	},
}

func init() {
	MainCmd.AddCommand(ioCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ioCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ioCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	ioCmd.PersistentFlags().StringP("top-words-by-google", "i", "/test", "dir for listening changes")
	ioCmd.PersistentFlags().StringP("piplineName", "p", "/test", "dir for listening changes")
	ioCmd.PersistentFlags().BoolP("show-top20k", "s", false, "dir for listening changes")

}
