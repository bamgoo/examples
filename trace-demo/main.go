package main

import (
	"fmt"
	"time"

	"github.com/bamgoo/bamgoo"
	. "github.com/bamgoo/base"
	_ "github.com/bamgoo/builtin"
	"github.com/bamgoo/http"
	_ "github.com/bamgoo/trace"
	_ "github.com/bamgoo/trace-file"
	_ "github.com/bamgoo/trace-greptime"
)

func main() {
	bamgoo.Go()
}

func init() {

	bamgoo.Register("index", http.Router{
		Uri: "/", Name: "index", Desc: "index",
		Action: func(ctx *http.Context) {
			ctx.Text("hello bamgoo.")
		},
	})

	bamgoo.Register("trace.child", bamgoo.Service{
		Name: "子调用", Desc: "trace child service",
		Action: func(ctx *bamgoo.Context) Map {
			ctx.Trace("搞飞机了这里")
			time.Sleep(10 * time.Millisecond)
			return Map{"ok": true, "at": time.Now().UnixMilli()}
		},
	})

	bamgoo.Register(bamgoo.START, bamgoo.Trigger{
		Name: "Trace Demo",
		Desc: "emit trace spans on startup",
		Action: func(ctx *bamgoo.Context) {
			data := ctx.Invoke("trace.child", Map{"from": "startup"})
			if res := ctx.Result(); res != nil && res.Fail() {

				return
			}

			span := ctx.Begin("开始")
			span.End()

			fmt.Println("trace demo done", data)
		},
	})
}
