package config

import (
	"github.com/gofsd/fsd/pkg/util"

	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	cfg = util.Config
	//ConfCmd conf cmd
	ConfCmd = &cobra.Command{
		Use:   "conf",
		Short: "add, update, delete and get config for command",
		Long:  `Interface to setup config for all cmds`,
		Run: func(cmd *cobra.Command, args []string) {
			confPathJoinedByDots := strings.Join(strings.Split(cmd.CommandPath(), " "), ".")
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				val, _ := cmd.Flags().GetString("create")
				print("from all: ", flag.Name, flag.Value.Type(), val, "\n")
				cfg.Set(confPathJoinedByDots, "test")
			})
			print(confPathJoinedByDots)
			//vi.Set()
		},
	}
)

func init() {
	ConfCmd.Flags().StringP("create", "c", "./input", "dir for search json files to convert to json")
	ConfCmd.Flags().StringP("read", "r", "./output", "dir for create go files with same folder struct like input dir")
	ConfCmd.Flags().StringP("update", "u", "./output", "dir for create go files with same folder struct like input dir")
	ConfCmd.Flags().StringP("delete", "d", "./output", "dir for create go files with same folder struct like input dir")
}
