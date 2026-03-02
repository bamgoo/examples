package main

import (
	"fmt"
	"time"

	_ "github.com/infrago/builtin"

	. "github.com/infrago/base"
	"github.com/infrago/infra"
	"github.com/infrago/web"
)

func main() {
	infra.Go()
}

func init() {
	infra.Register("*.access", web.Filter{
		Name: "请求日志",
		Serve: func(ctx *web.Context) {
			start := time.Now()
			fmt.Printf("[web] request site=%s host=%s domain=%s root=%s method=%s path=%s route=%s\n",
				ctx.Site, ctx.Host, ctx.Domain, ctx.RootDomain, ctx.Method, ctx.Path, ctx.Name)
			ctx.Next()
			fmt.Printf("[web] response host=%s method=%s path=%s status=%d cost=%s\n", ctx.Host, ctx.Method, ctx.Path, ctx.Code, time.Since(start))
		},
	})

	infra.Register("*.home", web.Router{
		Uri:  "/",
		Name: "通用首页",
		Action: func(ctx *web.Context) {
			ctx.Text(fmt.Sprintf("site=%s route=%s host=%s domain=%s root=%s path=%s",
				ctx.Site, ctx.Name, ctx.Host, ctx.Domain, ctx.RootDomain, ctx.Path))
		},
	})

	infra.Register("*.custom", web.Handler{
		Name: "统一错误处理",
		NotFound: func(ctx *web.Context) {
			ctx.JSON(Map{
				"code": 404,
				"msg":  "route not found",
				"site": ctx.Site,
				"host": ctx.Host,
				"path": ctx.Path,
			}, 404)
		},
		Error: func(ctx *web.Context) {
			msg := "internal server error"
			if res := ctx.Result(); res != nil && res.Error() != "" {
				msg = res.Error()
			}
			ctx.JSON(Map{
				"code": 500,
				"msg":  msg,
				"site": ctx.Site,
				"host": ctx.Host,
				"path": ctx.Path,
			}, 500)
		},
	})

	// 空站点（default）：只有以 "." 开头注册的才会进入这里。
	infra.Register(".index", web.Router{
		Uri:  "/",
		Name: "空站点首页",
		Action: func(ctx *web.Context) {
			ctx.JSON(Map{
				"site":       ctx.Site,
				"route":      ctx.Name,
				"url":        ctx.Url.Routo("www.index"),
				"host":       ctx.Host,
				"domain":     ctx.Domain,
				"rootDomain": ctx.RootDomain,
				"tip":        "matched empty-site router",
			})
		},
	})

	infra.Register("www.index", web.Router{
		Uri:  "/",
		Name: "WWW 首页",
		Action: func(ctx *web.Context) {
			ctx.Text("welcome to www site")
		},
	})

	infra.Register("www.about", web.Router{
		Uri:  "/about",
		Name: "WWW About",
		Action: func(ctx *web.Context) {
			ctx.Text("www about")
		},
	})

	infra.Register("file.get", web.Router{
		Uri:  "/download/{name}",
		Name: "文件下载模拟",
		Action: func(ctx *web.Context) {
			ctx.JSON(Map{
				"site": "file",
				"name": ctx.Params["name"],
				"url":  ctx.Path,
			})
		},
	})

	infra.Register("file.list", web.Router{
		Uri:  "/list",
		Name: "文件列表模拟",
		Action: func(ctx *web.Context) {
			ctx.JSON(Map{
				"site":  "file",
				"files": []string{"a.jpg", "b.pdf", "c.mp4"},
			})
		},
	})

	infra.Register("api.ping", web.Router{
		Uri:  "/ping",
		Name: "Ping",
		Action: func(ctx *web.Context) {
			ctx.JSON(Map{
				"ok":         true,
				"site":       "api",
				"domain":     ctx.Domain,
				"rootDomain": ctx.RootDomain,
			})
		},
	})

	infra.Register("api.user", web.Router{
		Uri:  "/user/{id}",
		Name: "用户详情",
		Action: func(ctx *web.Context) {
			ctx.JSON(Map{
				"site": "api",
				"id":   ctx.Params["id"],
			})
		},
	})

	infra.Register("api.fail", web.Router{
		Uri:  "/fail",
		Name: "错误模拟",
		Action: func(ctx *web.Context) {
			ctx.Error(infra.Fail.With("api fail test"))
		},
	})
}
