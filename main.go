package main

import (
	"github.com/gofsd/fsd/cmd"
	"github.com/gofsd/fsd/pkg/util"
)

func main() {
	util.HandleError(cmd.Execute())
}
