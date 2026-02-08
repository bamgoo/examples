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
	bamgoo.Register(bamgoo.START, bamgoo.Trigger{
		Name: "启动", Desc: "启动",
		Action: func(ctx *bamgoo.Context) {
			data := ctx.Invoke("test.get", Map{"msg": "msg from examples."})
			res := ctx.Result()

			fmt.Println("ssss", res, data)

			_, err := mutex.Lock("test", time.Minute)
			fmt.Println("lock", err)
		},
	})

}
