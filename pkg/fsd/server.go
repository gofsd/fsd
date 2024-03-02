package fsd

import (
	"bytes"
	"fmt"
	"sync"

	"time"

	"github.com/gin-gonic/gin"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin/binding"
	"github.com/gofsd/fsd/pkg/log"
	"github.com/spf13/cobra"
)

var r *gin.Engine

type SafeStdOut struct {
	mu   sync.Mutex
	b, e *bytes.Buffer
}

var std SafeStdOut

func getGinEngin(rootCmd *cobra.Command) (r *gin.Engine) {
	binding.Validator = new(defaultValidator)
	std.b = bytes.NewBuffer([]byte{})
	std.e = bytes.NewBuffer([]byte{})
	r = gin.New()

	logger := log.Get("")

	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	r.Use(ginzap.RecoveryWithZap(logger, true))
	r.POST("exec", Execute(rootCmd))
	r.GET("version", Version(rootCmd))
	r.GET("test", Test(rootCmd))

	return
}

func Start(port string, rootCmd *cobra.Command) {
	r = getGinEngin(rootCmd)
	r.Run(fmt.Sprintf(":%s", port))
}
