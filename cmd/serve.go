/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/gofsd/fsd/pkg/server"
	"github.com/spf13/cobra"
)

var port *string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")

		server.Start(config.ConfigPath, config.CmdPath, *port)
	},
}

func init() {
	MainCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//serveCmd.PersistentFlags().String("port", "", "App port")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//serveCmd.Flags().StringVar(port, "port", "8081", "App port")
}
