package main

import (
	"fmt"
	"time"

	. "github.com/bamgoo/base"
	_ "github.com/bamgoo/builtin"
	"github.com/bamgoo/mutex"

	"github.com/bamgoo/bamgoo"
	"github.com/bamgoo/http"
	"github.com/bamgoo/log"
)

func main() {
	bamgoo.Go()
}

func init() {

	bamgoo.Register("index", http.Router{
		Uri: "/", Name: "首页", Desc: "首页",
		Action: func(ctx *http.Context) {
			log.Debug("what")

			dri, err := bamgoo.Use[MailProvider]("sendcloud", Map{"msg": "setting"})
			log.Debug("use", err)
			if err == nil {
				data, res := dri.Send(Map{})
				log.Debug("send", data, res)
			}

			ctx.Text("hello world.")
		},
	})

	bamgoo.Register("www.index", http.Router{
		Uri: "/test", Name: "首页", Desc: "首页",
		Action: func(ctx *http.Context) {
			ctx.Text("hello test world.")
		},
	})

	// bamgoo.Register("www.index", http.Router{
	// 	Uri: "/", Name: "首页", Desc: "首页",
	// 	Action: func(ctx *http.Context) {
	// 		cache.Write("key", Map{"msg": "msg from cache."}, time.Second*10)
	// 		ctx.Text("hello world.")
	// 	},
	// })
	// bamgoo.Register("www.json", http.Router{
	// 	Uri: "/json", Name: "JSON", Desc: "JSON",
	// 	Action: func(ctx *http.Context) {
	// 		data, _ := cache.Read("key")
	// 		ctx.Echo(nil, Map{
	// 			"msg":   "hello world.",
	// 			"cache": data,
	// 		})
	// 	},
	// })

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
