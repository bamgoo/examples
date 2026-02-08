package main

import (
	"fmt"
	"time"

	. "github.com/bamgoo/base"
	_ "github.com/bamgoo/builtin"
	"github.com/bamgoo/mutex"

	"github.com/bamgoo/bamgoo"
)

func main() {
	bamgoo.Go()
}

func init() {

	bamgoo.Register("test.get", bamgoo.Service{
		Name: "查询", Desc: "查询",
		Action: func(ctx *bamgoo.Context) (Map, Res) {
			return Map{"msg": "retry from node 1"}, bamgoo.Retry
		},
	})

	bamgoo.Register(bamgoo.START, bamgoo.Trigger{
		Action: func(ctx *bamgoo.Context) {

			_, err := mutex.Lock("test", time.Minute)
			fmt.Println("lock", err)

			fmt.Println("start....")
		},
	})

}
