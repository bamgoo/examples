package main

import (
	"fmt"

	_ "github.com/bamgoo/builtin"

	"github.com/bamgoo/bamgoo"
	"github.com/bamgoo/web"
)

func main() {
	bamgoo.Go()
}

func init() {
	bamgoo.Register("*.logger", web.Filter{
		Name: "请求日志",
		Request: func(ctx *web.Context) {
			fmt.Printf("[web] host=%s method=%s path=%s route=%s\n", ctx.Host, ctx.Method, ctx.Path, ctx.Name)
			ctx.Next()
		},
	})

	bamgoo.Register("*.index", web.Router{
		Uri:  "/",
		Name: "站点首页",
		Action: func(ctx *web.Context) {
			ctx.Text(fmt.Sprintf("host=%s path=%s", ctx.Host, ctx.Path))
		},
	})

	bamgoo.Register("www.about", web.Router{
		Uri:  "/about",
		Name: "WWW About",
		Action: func(ctx *web.Context) {
			ctx.Text("www about")
		},
	})

	bamgoo.Register("user.profile", web.Router{
		Uri:  "/profile/{id}",
		Name: "用户资料",
		Action: func(ctx *web.Context) {
			ctx.Text(fmt.Sprintf("user profile id=%v host=%s", ctx.Params["id"], ctx.Host))
		},
	})

	bamgoo.Register("sys.health", web.Router{
		Uri:  "/health",
		Name: "系统健康检查",
		Action: func(ctx *web.Context) {
			ctx.JSON(map[string]any{"ok": true, "site": "sys", "host": ctx.Host})
		},
	})
}
