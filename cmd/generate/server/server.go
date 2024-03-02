package server

import (
	"github.com/gofsd/fsd/cmd/generate/template"
	"github.com/gofsd/fsd/pkg/util"

	"github.com/spf13/cobra"
)

var (
	rootCmdPath string
	cfg         = util.GetCfg()
	//GenerateCmd cmd
	ServerCmd = &cobra.Command{
		Use:   "server",
		Short: "generate server for fsd cli",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {

		},
		Run: func(cmd *cobra.Command, args []string) {
			template.Generate()
		},
	}
)

func init() {
	ServerCmd.PersistentFlags().StringVarP(&rootCmdPath, "config", "c", "config file (default is $HOME/.cobra.yaml)", "")
	ServerCmd.Flags().StringP("use", "u", "newCmd", "cmd name (required)")
}
