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
			return Map{"msg": "retry from node3"}, bamgoo.OK
		},
	})

	bamgoo.Register("test.node3", bamgoo.Service{
		Name: "test", Desc: "test",
		Action: func(ctx *bamgoo.Context) (Map, Res) {
			return Map{"msg": "retry from node3"}, bamgoo.OK
		},
	})

	bamgoo.Register(bamgoo.START, bamgoo.Trigger{
		Action: func(ctx *bamgoo.Context) {
			fmt.Println("node3 start....")
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

}
