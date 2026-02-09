package main

import (
	"fmt"
	"time"

	. "github.com/bamgoo/base"
	_ "github.com/bamgoo/builtin"
	"github.com/bamgoo/cache"
	"github.com/bamgoo/mutex"

	"github.com/bamgoo/bamgoo"
	"github.com/bamgoo/web"
)

func main() {
	bamgoo.Go()
}

func init() {
	bamgoo.Register("sys.index", web.Router{
		Uri: "/", Name: "首页", Desc: "首页",
		Action: func(ctx *web.Context) {
			ctx.Text("hello sys world.")
		},
	})

	bamgoo.Register("www.index", web.Router{
		Uri: "/", Name: "首页", Desc: "首页",
		Action: func(ctx *web.Context) {
			cache.Write("key", Map{"msg": "msg from cache."}, time.Second*10)
			ctx.Text("hello world.")
		},
	})
	bamgoo.Register("www.json", web.Router{
		Uri: "/json", Name: "JSON", Desc: "JSON",
		Action: func(ctx *web.Context) {
			data, _ := cache.Read("key")
			ctx.Echo(nil, Map{
				"msg":   "hello world.",
				"cache": data,
			})
		},
	})

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
