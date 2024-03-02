package macro

import (
	"path/filepath"

	"github.com/gofsd/fsd/pkg/util"

	"os"

	"github.com/spf13/cobra"
)

const (
	//MacrosFolder folder with all macros
	macrosFolder = "macros"
)

var (
	rootPath   = util.FindRootDir(util.CliName + "." + util.ConfigType)
	macrosPath = filepath.Join(rootPath, macrosFolder)
	macro      = make([]string, 0)
	cfg        = util.Config
	//MacroCmd cmd
	MacroCmd = &cobra.Command{
		Use:   "macro",
		Short: "macro for fsd cli",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			cfg.WriteCfgToFile(cmd)
		},
	}
)

func init() {
	MacroCmd.PersistentFlags().StringP("use", "u", "default", "macro to use")
	cfg.SetCmd(MacroCmd).Bind("use")

	MacroCmd.AddCommand(execCmd)
	MacroCmd.AddCommand(addCmd)
	err := os.MkdirAll(macrosPath, 0777)
	util.HandleError(err)
}
