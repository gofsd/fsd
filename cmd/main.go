/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	//add
	"github.com/gofsd/fsd/pkg/cfg"
	"github.com/gofsd/fsd/pkg/prompt"

	//add
	"github.com/spf13/cobra"
)

const (
	CREATE uint8 = iota + 1
	READ
	UPDATE
	DELETE
)

const (
	SERVER uint8 = iota + 1
	CLIENT
)

// add
var (
	config = cfg.GetCfg()
	DbName = "fsd.bolt"
	Port   = "32104"
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
		config.SetCommandPath(cmd.CommandPath())
		prompt.New().Run()
		//add
	},
}

func init() {
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
