package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Dialog(rootCmd *cobra.Command) *cobra.Command {
	var (
		Name   string
		Action uint8
		ID     uint32
	)
	var dialogCmd = &cobra.Command{
		Use:   "dialog",
		Short: "Create dialog in telegram",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("dialog called")
		},
	}
	dialogCmd.Flags().Uint8VarP(&Action, "action", "a", Action, "Actions: create is 1, read is 2, update is 3, delete is 4")
	dialogCmd.Flags().StringVarP(&Name, "name", "n", Name, "Dialog Name")
	dialogCmd.Flags().Uint32VarP(&ID, "id", "i", ID, "Dialog ID")

	rootCmd.AddCommand(dialogCmd)

	return dialogCmd
}

func CreateDialog(rootCmd *cobra.Command) *cobra.Command {
	var (
		Q, A, D string
	)
	var dialogCmd = &cobra.Command{
		Use:   "create",
		Short: "Create dialog",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("dialog called")
		},
	}

	dialogCmd.Flags().StringVarP(&D, "dialog", "d", D, "Dialog id")
	dialogCmd.Flags().StringVarP(&Q, "question", "q", Q, "Question")
	dialogCmd.Flags().StringVarP(&A, "answer", "a", A, "Answer")

	rootCmd.AddCommand(dialogCmd)

	return dialogCmd
}

func DeleteDialog(rootCmd *cobra.Command) *cobra.Command {
	var (
		Q, A, D string
	)
	var dialogCmd = &cobra.Command{
		Use:   "create",
		Short: "Create dialog",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("dialog called")
		},
	}

	dialogCmd.Flags().StringVarP(&D, "dialog", "d", D, "Dialog id")
	dialogCmd.Flags().StringVarP(&Q, "question", "q", Q, "Question")
	dialogCmd.Flags().StringVarP(&A, "answer", "a", A, "Answer")

	rootCmd.AddCommand(dialogCmd)

	return dialogCmd
}

func UpdateDialog(rootCmd *cobra.Command) *cobra.Command {
	var (
		Q, A, D string
	)
	var dialogCmd = &cobra.Command{
		Use:   "create",
		Short: "Create dialog",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("dialog called")
		},
	}

	dialogCmd.Flags().StringVarP(&D, "dialog", "d", D, "Dialog id")
	dialogCmd.Flags().StringVarP(&Q, "question", "q", Q, "Question")
	dialogCmd.Flags().StringVarP(&A, "answer", "a", A, "Answer")

	rootCmd.AddCommand(dialogCmd)

	return dialogCmd
}

func ViewDialog(rootCmd *cobra.Command) *cobra.Command {
	var (
		Q, A, D string
	)
	var dialogCmd = &cobra.Command{
		Use:   "create",
		Short: "Create dialog",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("dialog called")
		},
	}

	dialogCmd.Flags().StringVarP(&D, "dialog", "d", D, "Dialog id")
	dialogCmd.Flags().StringVarP(&Q, "question", "q", Q, "Question")
	dialogCmd.Flags().StringVarP(&A, "answer", "a", A, "Answer")

	rootCmd.AddCommand(dialogCmd)

	return dialogCmd
}

func init() {
	DialogCmd := Dialog(MainCmd)
	CreateDialog(DialogCmd)
	DeleteDialog(DialogCmd)
	UpdateDialog(DialogCmd)
	ViewDialog(DialogCmd)

}
