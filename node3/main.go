package main

import (
	. "github.com/bamgoo/base"
	_ "github.com/bamgoo/builtin"

	"github.com/bamgoo/bamgoo"
)

func main() {
	bamgoo.Go()
}

func init() {

	bamgoo.Register("test.get", bamgoo.Service{
		Name: "查询", Desc: "查询",
		Action: func(ctx *bamgoo.Context) (Map, Res) {

			return Map{"msg": "fail from node 3"}, bamgoo.Fail
		},
	})

}
