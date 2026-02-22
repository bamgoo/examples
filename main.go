package main

import (
	"fmt"
	"time"

	. "github.com/bamgoo/base"
	_ "github.com/bamgoo/builtin"

	"github.com/bamgoo/cron"
	"github.com/bamgoo/log"
	"github.com/bamgoo/mutex"

	"github.com/bamgoo/bamgoo"
	"github.com/bamgoo/http"

	_ "github.com/bamgoo/cron-pgsql"
)

func main() {
	bamgoo.Go()
}

func init() {

	bamgoo.Register("cron.test", bamgoo.Method{
		Name: "test", Desc: "test",
		Action: func(ctx *bamgoo.Context) (Map, Res) {
			log.Debug("cron.test", time.Now())
			return nil, bamgoo.OK
		},
	})

	bamgoo.Register("test", cron.Job{
		Schedule: "*/10 * * * * *", Target: "cron.test",
	})

	bamgoo.Register("index", http.Router{
		Uri: "/", Name: "首页", Desc: "首页",
		Action: func(ctx *http.Context) {
			jobs := cron.GetJobs()
			count, logs := cron.GetLogs("test", 0, 10)
			ctx.JSON(Map{
				"count": count, "logs": logs,
				"jobs": jobs,
			})
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
