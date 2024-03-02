package macro

import (
	"encoding/json"
	"io/ioutil"
	"os/exec"

	"github.com/gofsd/fsd/pkg/util"

	"os"
	"path/filepath"

	"strings"

	"github.com/spf13/cobra"
)

var (
	execCmd = &cobra.Command{
		Use:   "exec",
		Short: "exec macros",

		Run: func(cmd *cobra.Command, args []string) {
			print("from exec")
			filePath := filepath.Join(macrosPath, cfg.GetString("use"))
			file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777)
			defer file.Close()
			macroBytes, err := ioutil.ReadAll(file)
			util.HandleError(err)
			err = json.Unmarshal([]byte(macroBytes), &macro)
			for i := 0; len(macro) > i; i++ {
				commandArr := strings.Split(macro[i], " ")
				c := exec.Command(commandArr[0], commandArr[1:]...)
				c.Stdin = os.Stdin
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				err = c.Run()
				util.HandleError(err)
			}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			cfg.WriteCfgToFile(cmd)
		},
	}
)

func init() {
}
