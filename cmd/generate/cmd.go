package generate

import (
	"github.com/gofsd/fsd/cmd/generate/server"
	"github.com/gofsd/fsd/cmd/generate/template"
	"github.com/gofsd/fsd/pkg/util"

	"github.com/spf13/cobra"
)

var (
	rootCmdPath string
	cfg         = util.GetCfg()
	//GenerateCmd cmd
	GenerateCmd = &cobra.Command{
		Use:   "generate",
		Short: "generate command for fsd cli",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {

		},
		Run: func(cmd *cobra.Command, args []string) {
			template.Generate()
		},
	}
)

func init() {
	GenerateCmd.AddCommand(server.ServerCmd)
	GenerateCmd.PersistentFlags().StringVarP(&rootCmdPath, "config", "c", "config file (default is $HOME/.cobra.yaml)", "")
	GenerateCmd.Flags().StringP("use", "u", "newCmd", "cmd name (required)")
}
