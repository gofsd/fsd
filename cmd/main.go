/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	//add

	types "github.com/gofsd/fsd-types"
	"github.com/gofsd/fsd/pkg/prompt"

	//add
	"github.com/spf13/cobra"
)

//add

// MainCmd represents the root command
var MainCmd = &cobra.Command{
	Use:   "fsd",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("root called")

		//add
		//config.SetCommandPath(cmd.CommandPath())
		prompt.New().Run()
		//add
	},
}

func init() {
	os.MkdirAll(types.DBFullName, 0666)
	MainCmd.SetOut(os.Stdout)
	MainCmd.SetErr(os.Stderr)

	cobra.OnInitialize(initConfig)

	MainCmd.PersistentFlags().String("foo", "", "A help for foo")

	MainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Execute() error {
	return MainCmd.Execute()
}

func initConfig() {
	//config.ReadCfgFromFile(nil)
	//config.SetMainCmd(MainCmd)
	//config.SetLog()

}
