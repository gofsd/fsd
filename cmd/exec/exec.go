package exec

import (
	"github.com/gofsd/fsd/cmd/generate/template"
	"github.com/gofsd/fsd/pkg/util"

	"github.com/spf13/cobra"
)

var (
	macroCmdPath string
	cfg          = util.GetCfg()
	//ExecCmd cmd
	ExecCmd = &cobra.Command{
		Use:   "macro",
		Short: "macro for fsd cli",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {

		},
		Run: func(cmd *cobra.Command, args []string) {
			template.Generate()
		},
	}
)

func init() {
	ExecCmd.PersistentFlags().StringVarP(&macroCmdPath, "config", "c", "config file (default is $HOME/.cobra.yaml)", "")
	ExecCmd.Flags().StringP("use", "u", "newCmd", "cmd name (required)")
}
