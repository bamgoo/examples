package main

import (
	"fmt"
	"time"

	"github.com/bamgoo/bamgoo"
	. "github.com/bamgoo/base"
	_ "github.com/bamgoo/builtin"
	_ "github.com/bamgoo/trace"
	_ "github.com/bamgoo/trace-file"
	_ "github.com/bamgoo/trace-greptime"
)

func main() {
	bamgoo.Go()
}

func init() {
	bamgoo.Register("trace.child", bamgoo.Service{
		Name: "子调用",
		Desc: "trace child service",
		Action: func(ctx *bamgoo.Context) (Map, Res) {
			ctx.Trace("trace.child.step", Map{"status": "ok", "step": "compute"})
			time.Sleep(10 * time.Millisecond)
			return Map{"ok": true, "at": time.Now().UnixMilli()}, bamgoo.OK
		},
	})

	bamgoo.Register(bamgoo.START, bamgoo.Trigger{
		Name: "Trace Demo",
		Desc: "emit trace spans on startup",
		Action: func(ctx *bamgoo.Context) {
			ctx.Trace("trace.start.before_invoke", Map{"status": "ok"})
			data := ctx.Invoke("trace.child", Map{"from": "startup"})
			if res := ctx.Result(); res != nil && res.Fail() {
				ctx.Trace("trace.start.error", Map{"status": "error", "error": res.Error()})
				return
			}
			ctx.Trace("trace.start.after_invoke", Map{"status": "ok", "data": data})
			fmt.Println("trace demo done", data)
		},
	})
}
