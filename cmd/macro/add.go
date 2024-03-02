package macro

import (
	"encoding/json"
	"io/ioutil"

	"os"
	"path/filepath"

	"github.com/gofsd/fsd/pkg/util"
	"github.com/spf13/cobra"
)

const (
// MacrosFolder folder with all macros
)

var (

	//MacroCmd cmd
	addCmd = &cobra.Command{
		Use:   "add",
		Short: "add command to macros",
		Run: func(cmd *cobra.Command, args []string) {
			filePath := filepath.Join(macrosPath, cfg.GetString("use"))
			file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777)
			defer file.Close()
			macroBytes, err := ioutil.ReadAll(file)
			util.HandleError(err)
			err = json.Unmarshal([]byte(macroBytes), &macro)
			for i := 0; len(args) > i; i++ {
				macro = append(macro, args[i])
			}
			macroBytes, err = json.Marshal(macro)
			util.HandleError(err)
			file.Truncate(0)
			file.WriteAt(macroBytes, 0)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			cfg.WriteCfgToFile(cmd)
		},
	}
)

func init() {
}
