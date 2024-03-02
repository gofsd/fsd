/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.SetCommandPath(cmd.CommandPath())
		config.Bind("url")
		config.Bind("id")
		config.Bind("uuid")
	},
	Run: func(cmd *cobra.Command, args []string) {
		var cmdParts []string
		cmdParts = append(cmdParts, "curl")
		cmdParts = append(cmdParts, "-H")
		cmdParts = append(cmdParts, "ID: "+config.GetString("usr.id"))

		cmdParts = append(cmdParts, "-H")
		cmdParts = append(cmdParts, "UUID: "+config.GetString("usr.uuid"))
		cmdParts = append(cmdParts, "-o")
		cmdParts = append(cmdParts, config.ConfigPath+"/"+config.GetString("backend_path")+"my.db")

		cmdParts = append(cmdParts, config.GetString("backup.url"))
		fmt.Printf("%+v", cmdParts)
		process := exec.Command(cmdParts[0], cmdParts[1:]...)
		print(process.String())
		process.Stdout = os.Stdout
		process.Run()
	},
}

func init() {
	MainCmd.AddCommand(backupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// backupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	backupCmd.PersistentFlags().StringP("url", "u", "./input", "backup url")
	backupCmd.PersistentFlags().StringP("id", "i", "./output", "id")
	backupCmd.PersistentFlags().StringP("uid", "d", "./output", "uid")

}
