/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	fsd "github.com/gofsd/fsd-types"
	c "github.com/gofsd/fsd/pkg/cmd"
	kv "github.com/gofsd/fsd/pkg/kv"
	s "github.com/gofsd/fsd/pkg/store"
	"github.com/spf13/cobra"
)

func Tree(rootCmd *cobra.Command) *cobra.Command {
	var (
		value, dbName           string
		action                  uint8
		id, defaultID, bucketID uint16
	)
	var treeCmd = &cobra.Command{
		Use:   "tree",
		Short: "Create tree",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			db := s.New(
				s.SetFullDbName(dbName),
				s.SetBucketName(kv.GetKeyFromInt(bucketID)),
			)
			pair := &fsd.Pair{}
			pair.SetID(uint64(id))
			pair.FromString(value)
			defer db.Close()
			h := c.Set(cmd).
				Equal(action, CREATE).
				AND().
				Equal(id, defaultID).
				HandleCRUD(db.JustCreate, pair).
				Equal(action, READ).
				AND().
				NotEqual(id, defaultID).
				HandleCRUD(db.JustRead, pair).
				Equal(action, UPDATE).
				AND().
				NotEqual(id, defaultID).
				HandleCRUD(db.JustUpdate, pair).
				Equal(action, DELETE).
				AND().
				NotEqual(id, defaultID).
				HandleCRUD(db.JustDelete, pair).
				Error().
				Equal(Output, "byte").
				HandleB(pair.Json).
				Equal(Output, "json").
				HandleB(pair.Json).
				Equal(Output, "string").
				JustStr(pair.String).
				Error()
			return h.E
		},
	}

	treeCmd.Flags().Uint8VarP(&action, "action", "a", action, "Actions: create is 1, read is 2, update is 3, delete is 4")
	treeCmd.MarkFlagRequired("action")
	treeCmd.Flags().StringVarP(&value, "value", "n", value, "Set value")
	treeCmd.MarkFlagRequired("value")
	treeCmd.Flags().Uint16VarP(&id, "id", "i", id, "Tree value id(use 0 if you use create action)")
	treeCmd.MarkFlagRequired("id")
	treeCmd.Flags().Uint16VarP(&bucketID, "bucket", "b", bucketID, "Bucket ID")
	treeCmd.Flags().StringVarP(&dbName, "db-name", "d", DbName, "Db Name")
	treeCmd.Flags().StringVarP(&Output, "output", "o", Output, "Parent ID")

	rootCmd.AddCommand(treeCmd)

	return treeCmd
}

func TreeLeaf(rootCmd *cobra.Command) *cobra.Command {
	var (
		Name, Parent string
		Action       uint8
		ID           uint32
	)
	var treeLeafCmd = &cobra.Command{
		Use:   "leaf",
		Short: "Create tree",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("dialog called")
		},
	}

	treeLeafCmd.Flags().Uint32VarP(&ID, "id", "i", ID, "Tree ID")
	treeLeafCmd.Flags().Uint8VarP(&Action, "action", "a", Action, "Actions: create is 1, read is 2, update is 3, delete is 4")
	treeLeafCmd.Flags().StringVarP(&Name, "name", "n", Name, "Dialog Name")
	treeLeafCmd.Flags().StringVarP(&Parent, "parent", "p", Parent, "Parent Name")

	rootCmd.AddCommand(treeLeafCmd)

	return treeLeafCmd
}

func TreeBranch(rootCmd *cobra.Command) *cobra.Command {
	var (
		Name, Parent string
		Action       uint8
		ID           uint32
	)
	var treeBranchCmd = &cobra.Command{
		Use:   "tree",
		Short: "Create tree",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("dialog called")
		},
	}

	treeBranchCmd.Flags().Uint32VarP(&ID, "id", "i", ID, "Tree ID")
	treeBranchCmd.Flags().Uint8VarP(&Action, "action", "a", Action, "Actions: create is 1, read is 2, update is 3, delete is 4")
	treeBranchCmd.Flags().StringVarP(&Name, "name", "n", Name, "Dialog Name")
	treeBranchCmd.Flags().StringVarP(&Parent, "parent", "p", Parent, "Parent Name")

	rootCmd.AddCommand(treeBranchCmd)

	return treeBranchCmd
}

func TreeCrown(rootCmd *cobra.Command) *cobra.Command {
	var (
		Name, Parent string
		Action       uint8
		ID           uint32
	)
	var treeCrownCmd = &cobra.Command{
		Use:   "tree",
		Short: "Create tree",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("dialog called")
		},
	}

	treeCrownCmd.Flags().Uint32VarP(&ID, "id", "i", ID, "Tree ID")
	treeCrownCmd.Flags().Uint8VarP(&Action, "action", "a", Action, "Actions: create is 1, read is 2, update is 3, delete is 4")
	treeCrownCmd.Flags().StringVarP(&Name, "name", "n", Name, "Dialog Name")
	treeCrownCmd.Flags().StringVarP(&Parent, "parent", "p", Parent, "Parent Name")

	rootCmd.AddCommand(treeCrownCmd)

	return treeCrownCmd
}

func TreeTrunk(rootCmd *cobra.Command) *cobra.Command {
	var (
		Name, Parent string
		Action       uint8
		ID           uint32
	)
	var treeTrunkCmd = &cobra.Command{
		Use:   "tree",
		Short: "Create tree",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("dialog called")
		},
	}

	treeTrunkCmd.Flags().Uint32VarP(&ID, "id", "i", ID, "Tree ID")
	treeTrunkCmd.Flags().Uint8VarP(&Action, "action", "a", Action, "Actions: create is 1, read is 2, update is 3, delete is 4")
	treeTrunkCmd.Flags().StringVarP(&Name, "name", "n", Name, "Dialog Name")
	treeTrunkCmd.Flags().StringVarP(&Parent, "parent", "p", Parent, "Parent Name")

	rootCmd.AddCommand(treeTrunkCmd)

	return treeTrunkCmd
}

func TreeTwig(rootCmd *cobra.Command) *cobra.Command {
	var (
		Name, Parent string
		Action       uint8
		ID           uint32
	)
	var treeTrunkCmd = &cobra.Command{
		Use:   "tree",
		Short: "Create tree",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("dialog called")
		},
	}

	treeTrunkCmd.Flags().Uint32VarP(&ID, "id", "i", ID, "Tree ID")
	treeTrunkCmd.Flags().Uint8VarP(&Action, "action", "a", Action, "Actions: create is 1, read is 2, update is 3, delete is 4")
	treeTrunkCmd.Flags().StringVarP(&Name, "name", "n", Name, "Dialog Name")
	treeTrunkCmd.Flags().StringVarP(&Parent, "parent", "p", Parent, "Parent Name")

	rootCmd.AddCommand(treeTrunkCmd)

	return treeTrunkCmd
}

func init() {
	TreeCmd := Tree(MainCmd)
	TreeLeaf(TreeCmd)
}
