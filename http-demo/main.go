package main

import (
	"fmt"
	"time"

	_ "github.com/infrago/builtin"
	"github.com/infrago/http"
	"github.com/infrago/infra"
)

func main() {
	infra.Go()
}

func init() {
	infra.Register("*.access", http.Filter{
		Name: "请求日志",
		Serve: func(ctx *http.Context) {
			start := time.Now()
			fmt.Printf("[http] request site=%s host=%s domain=%s root=%s method=%s path=%s route=%s\n",
				ctx.Site, ctx.Host, ctx.Domain, ctx.RootDomain, ctx.Method, ctx.Path, ctx.Name)
			ctx.Next()
			fmt.Printf("[http] response host=%s method=%s path=%s status=%d cost=%s\n", ctx.Host, ctx.Method, ctx.Path, ctx.Code, time.Since(start))
		},
	})

	infra.Register("*.home", http.Router{
		Uri:  "/",
		Name: "通用首页",
		Action: func(ctx *http.Context) {
			ctx.Text(fmt.Sprintf("site=%s route=%s host=%s domain=%s root=%s path=%s",
				ctx.Site, ctx.Name, ctx.Host, ctx.Domain, ctx.RootDomain, ctx.Path))
			ctx.View("index")
		},
	})

}
