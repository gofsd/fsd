package fsd

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	typs "github.com/gofsd/fsd-types"
	command "github.com/gofsd/fsd/pkg/cmd"
	"github.com/spf13/cobra"
)

func Execute(rootCmd *cobra.Command) func(*gin.Context) {
	cmdStore := command.Store()
	return func(ctx *gin.Context) {
		var response typs.CommandResponse

		fn := func() {
			defer Duration(time.Now(), &response)
			var (
				c       typs.Command
				mainCmd = rootCmd
				out     []byte
			)
			if err := ctx.ShouldBindJSON(&c); err != nil {
				var e typs.Error
				e.Error = err.Error()
				e.Command = &c
				ctx.JSON(http.StatusBadRequest, e)
				return
			} else {
				cmd, args, err := mainCmd.Find(c.Name)
				if err != nil {
					var e typs.Error
					e.Error = err.Error()
					e.Command = &c
					ctx.JSON(http.StatusBadRequest, e)
					return
				}
				for _, kv := range c.Flags {
					cmd.Flags().Set(kv.K, kv.V)
				}
				std.mu.Lock()
				defer std.mu.Unlock()

				cmd.SetOut(std.b)
				cmd.SetErr(std.e)

				cmd.RunE(cmd, args)
				out, _ = io.ReadAll(std.b)
				std.b.Reset()

				if len(out) == 0 {
					out, _ = io.ReadAll(std.e)
					std.e.Reset()
					var e typs.Error
					e.Error = string(out)
					e.Command = &c
					ctx.JSON(http.StatusBadRequest, e)
					return
				}
				response = cmdStore.Save(c, out)

			}
		}
		fn()
		ctx.JSON(200, response)
		return
	}
}

func Duration(invocation time.Time, resp *typs.CommandResponse) {
	resp.Duration = int64(time.Since(invocation))
}

func Version(rootCmd *cobra.Command) func(*gin.Context) {
	return func(ctx *gin.Context) {
		cmd, args, err := rootCmd.Find([]string{"version"})
		if err != nil {
			ctx.JSON(200, err)
		} else {
			std.mu.Lock()
			defer std.mu.Unlock()
			cmd.SetOut(std.b)
			cmd.RunE(cmd, args)
			out, _ := io.ReadAll(std.b)
			std.b.Reset()
			ctx.Data(200, "text/html", out)
		}
	}
}

func Test(rootCmd *cobra.Command) func(*gin.Context) {
	return func(ctx *gin.Context) {
		r, err := http.Get("https://google.com")
		if err != nil {
			ctx.JSON(200, err)
		} else {
			out, _ := io.ReadAll(r.Body)
			std.b.Reset()
			ctx.Data(200, "text/html", out)
		}
	}
}

func init() {
}
