package convert

import (
	"os"
	"strings"

	"github.com/gofsd/fsd/pkg/util"

	"github.com/ChimeraCoder/gojson"
	"github.com/spf13/cobra"
)

var (
	cfg = util.GetCfg()
	//JSON2Go cmd
	JSON2Go = &cobra.Command{
		Use:   "json2go",
		Short: "convert json to go structure",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cfg.Bind("inputDir")
			cfg.Bind("outputDir")

		},
		PreRun: func(cmd *cobra.Command, args []string) {

		},
		Run: func(cmd *cobra.Command, args []string) {
			var (
				res  []byte
				err  error
				file *os.File
			)
			os.MkdirAll(cfg.GetString("inputDir"), 0777)
			os.MkdirAll(cfg.GetString("outputDir"), 0777)

			//		print("viper data", viper.GetString("inputDir"))
			util.ExecForFilesWithExtension(cfg.GetString("inputDir"), ".json", func(fullName string, name string) {
				print("fullName: ", fullName)
				file, err = os.Open(fullName)
				util.HandleError(err)
				res, err = gojson.Generate(file, gojson.ParseJson, "someStruct", "somepkg", []string{"json"}, true, true)
				util.HandleError(err)
				name = strings.Replace(name, ".json", ".go", -1)
				file, err = os.OpenFile(cfg.GetString("outputDir")+"/"+name, os.O_RDWR|os.O_CREATE, 0666)
				util.HandleError(err)
				err = file.Truncate(0)
				util.HandleError(err)
				_, err = file.Write(res)
				util.HandleError(err)
			})
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			cfg.WriteCfgToFile(cmd)
		},
	}
)

func init() {
	JSON2Go.PersistentFlags().StringP("inputDir", "i", "./input", "dir for search json files to convert to json")
	JSON2Go.PersistentFlags().StringP("outputDir", "o", "./output", "dir for create go files with same folder struct like input dir")
}
