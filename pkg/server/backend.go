package server

import (
	server "github.com/gofsd/fsd/pkg/server/http"

	"os"
	"os/signal"
	"syscall"
)

func Start(rootPath, cmdPath, port string) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		print(sig.String())
		os.Exit(0)
	}()
	server.RootPath = rootPath
	server.CmdPath = cmdPath
	server.Create()
	server.Run(port)
}
