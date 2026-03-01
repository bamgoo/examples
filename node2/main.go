package main

import (
	"fmt"

	. "github.com/infrago/base"
	_ "github.com/infrago/builtin"
	_ "github.com/infrago/bus-redis"

	"github.com/infrago/infra"
	"github.com/infrago/bus"
	"github.com/infrago/http"
)

func main() {
	infra.Go()
}

func init() {

	infra.Register("test", infra.Service{
		Name: "测试服务", Desc: "测试服务",
		Action: func(ctx *infra.Context) (Map, Res) {
			return Map{"msg": "retry from node2"}, infra.OK
		},
	})

	infra.Register("test.node2", infra.Service{
		Name: "test", Desc: "test",
		Action: func(ctx *infra.Context) (Map, Res) {
			return Map{"msg": "retry from node2"}, infra.OK
		},
	})

	infra.Register(infra.START, infra.Trigger{
		Action: func(ctx *infra.Context) {
			fmt.Println("node2 start....")
		},
	})

	infra.Register("index", http.Router{
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
