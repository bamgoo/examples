package main

import (
	"fmt"

	. "github.com/bamgoo/base"
	_ "github.com/bamgoo/builtin"
	_ "github.com/bamgoo/bus-redis"

	"github.com/bamgoo/bamgoo"
	"github.com/bamgoo/bus"
	"github.com/bamgoo/http"
)

func main() {
	bamgoo.Go()
}

func init() {

	bamgoo.Register("test", bamgoo.Service{
		Name: "测试服务", Desc: "测试服务",
		Action: func(ctx *bamgoo.Context) (Map, Res) {
			return Map{"msg": "retry from node4"}, bamgoo.OK
		},
	})

	bamgoo.Register("test.node4", bamgoo.Service{
		Name: "test", Desc: "test",
		Action: func(ctx *bamgoo.Context) (Map, Res) {
			return Map{"msg": "retry from node4"}, bamgoo.OK
		},
	})

	bamgoo.Register(bamgoo.START, bamgoo.Trigger{
		Action: func(ctx *bamgoo.Context) {
			fmt.Println("node4 start....")
		},
	})

	bamgoo.Register("index", http.Router{
		Uri: "/", Name: "首页", Desc: "首页",
		Action: func(ctx *http.Context) {

			nodes := bus.ListNodes()
			services := bus.ListServices()

			ctx.JSON(Map{
				"nodes": nodes, "services": services,
			})
		},
	})

	bamgoo.Register("invoke", http.Router{
		Uri: "/invoke", Name: "invoke", Desc: "invoke",
		Action: func(ctx *http.Context) {

			data := ctx.Invoke("test.node3")

			ctx.JSON(Map{
				"data": data, "result": ctx.Result(),
			})
		},
	})

}
