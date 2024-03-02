/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/gofsd/fsd/pkg/bolt"

	"github.com/spf13/cobra"
)

// addBoltDatabaseCmd represents the addBoltDatabase command
var addBoltDatabaseCmd = &cobra.Command{
	Use:   "add",
	Short: "Open or create new bolt db",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range args {
			db := bolt.CreateOrOpenDB(v)
			if db.IsOpen() {
				fmt.Printf("Path: %s; Page size: %d; Db size: %d\n", db.Path(), db.GetPageSize(), db.GetSize())
			}
			db.Close()
		}
	},
}

func init() {
	boltDatabaseCmd.AddCommand(addBoltDatabaseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addBoltDatabaseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addBoltDatabaseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
