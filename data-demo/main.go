package main

import (
	"fmt"
	"time"

	"github.com/infrago/infra"
	. "github.com/infrago/base"
	_ "github.com/infrago/builtin"
	"github.com/infrago/data"
	_ "github.com/infrago/data-mysql"
	_ "github.com/infrago/data-postgresql"
	_ "github.com/infrago/data-sqlite"
	"github.com/infrago/http"
	"github.com/infrago/log"
)

func main() {
	infra.Go()
}

func init() {
	infra.Register("user", data.Table{
		Name: "用户",
		Key:  "id",
		Fields: Vars{
			"id":          Var{Type: "int"},
			"name":        Var{Type: "string", Required: true},
			"status":      Var{Type: "string"},
			"login_times": Var{Type: "int"},
			"profile":     Var{Type: "jsonb"},
			"tags":        Var{Type: "strings"},
			"created":     Var{Type: "timestamp"},
		},
	})

	infra.Register("order", data.Table{
		Name: "订单",
		Key:  "id",
		Fields: Vars{
			"id":      Var{Type: "int"},
			"user_id": Var{Type: "int", Required: true},
			"amount":  Var{Type: "float"},
			"status":  Var{Type: "string"},
			"created": Var{Type: "timestamp"},
		},
	})
	infra.Register("order_agg", data.Table{
		Name: "订单聚合快照", Desc: "asdf",
		Schema: "public", Key: "user_id",
		Fields: Vars{
			"user_id": Var{Type: "int"},
			"total":   Var{Type: "float"},
			"cnt":     Var{Type: "int"},
			"updated": Var{Type: "timestamp"},
		},
	})

	infra.Register("user", data.Watcher{
		Action: func(m data.Mutation) {
			log.Debug("data watcher any", m.Table, m.Op, m.Rows)
		},
	})

	infra.Register("order", data.Watcher{
		Insert: func(m data.Mutation) {
			log.Debug("data watcher create", m.Table, m.Key)
		},
		Update: func(m data.Mutation) {
			log.Debug("data watcher update", m.Table, m.Rows)
		},
		Delete: func(m data.Mutation) {
			log.Debug("data watcher delete", m.Table, m.Rows)
		},
	})

	infra.Register(infra.START, infra.Trigger{
		Name: "Data Demo", Desc: "Run data module demo on startup",
		Action: func(ctx *infra.Context) {
			db := data.Base("main")
			defer db.Close()
			fmt.Println("capabilities", db.Capabilities())

			q, _ := data.Parse(Map{
				"status": "active",
				"$limit": 10,
			})
			log.Debug("q", q)

			adb := data.Base("analytics")
			defer adb.Close()

			_ = db.Exec(`CREATE TABLE IF NOT EXISTS "user" (
				"id" INTEGER PRIMARY KEY,
				"name" TEXT NOT NULL,
				"status" TEXT,
				"login_times" INTEGER DEFAULT 0,
				"profile" TEXT,
				"tags" TEXT,
				"created" DATETIME
			)`)

			_ = db.Exec(`CREATE TABLE IF NOT EXISTS "order" (
				"id" INTEGER PRIMARY KEY,
				"user_id" INTEGER NOT NULL,
				"amount" REAL,
				"status" TEXT,
				"created" DATETIME
			)`)
			_ = adb.Exec(`CREATE TABLE IF NOT EXISTS "order_agg" (
				"user_id" INTEGER PRIMARY KEY,
				"total" REAL,
				"cnt" INTEGER,
				"updated" DATETIME
			)`)

			user := db.Table("user").Upsert(Map{
				"$set": Map{
					"name":    "Alice",
					"status":  "active",
					"profile": Map{"city": "Shanghai", "vip": true},
					"tags":    []string{"go", "db", "cloud"},
					"created": time.Now(),
				},
				"$inc": Map{
					"login_times": ASC,
				},
			}, Map{"id": 1001})
			if db.Error() != nil {
				fmt.Println("upsert user error:", db.Error())
				return
			}

			_ = db.Table("order").Insert(Map{"id": 1, "user_id": 1001, "amount": 39.5, "status": "paid", "created": time.Now()})
			_ = db.Table("order").Insert(Map{"id": 2, "user_id": 1001, "amount": 72.0, "status": "paid", "created": time.Now()})
			_ = db.Table("order").Insert(Map{"id": 3, "user_id": 1001, "amount": 9.0, "status": "new", "created": time.Now()})

			rows := db.Table("order").Aggregate(Map{
				"status": "paid",
				"$group": []string{"user_id"},
				"$agg": Map{
					"total": Map{"sum": "amount"},
					"cnt":   Map{"count": "*"},
				},
				"$having": Map{
					"total": Map{"$gt": 50},
				},
			})
			if db.Error() != nil {
				fmt.Println("aggregate error:", db.Error())
				return
			}
			for _, row := range rows {
				_ = adb.Table("order_agg").Upsert(Map{
					"$set": Map{
						"total":   row["total"],
						"cnt":     row["cnt"],
						"updated": time.Now(),
					},
				}, Map{"user_id": row["user_id"]})
			}

			tagged := db.Table("user").Query(Map{
				"tags":  Map{"$contains": []string{"go"}},
				"$sort": Map{"id": ASC},
			})
			if db.Error() != nil {
				fmt.Println("contains query error:", db.Error())
			}

			_ = db.Table("order").ScanN(2, func(item Map) Res {
				fmt.Println("range item", item)
				return infra.OK
			}, Map{
				"$sort": Map{"id": ASC},
			})

			_, _ = db.View("asdf").Slice(0, 20, Map{})

			fmt.Println("user", user)
			fmt.Println("aggregate rows", rows)
			fmt.Println("tagged users", tagged)
		},
	})

	infra.Register("data.user", http.Router{
		Uri:  "/data/user/{id}",
		Name: "查询用户",
		Desc: "按ID查询用户",
		Action: func(ctx *http.Context) {
			db := data.Base("main")
			defer db.Close()

			user := db.Table("user").Entity(ctx.Params["id"])
			if db.Error() != nil {
				ctx.JSON(Map{"ok": false, "error": db.Error().Error()})
				return
			}

			ctx.JSON(Map{"ok": true, "user": user})
		},
	})

	infra.Register("data.orders", http.Router{
		Uri:  "/data/orders",
		Name: "订单聚合",
		Desc: "订单聚合统计",
		Action: func(ctx *http.Context) {
			db := data.Base("main")
			defer db.Close()

			rows := db.Table("order").Aggregate(Map{
				"status": "paid",
				"$group": []string{"user_id"},
				"$agg": Map{
					"total": Map{"sum": "amount"},
					"cnt":   Map{"count": "*"},
				},
			})
			if err := db.Error(); err != nil {
				ctx.JSON(Map{"ok": false, "error": err.Error()})
				return
			}

			ctx.JSON(Map{"ok": true, "items": rows})
		},
	})

	infra.Register("data.capabilities", http.Router{
		Uri:  "/data/capabilities",
		Name: "能力矩阵",
		Desc: "查看当前数据驱动能力",
		Action: func(ctx *http.Context) {
			mainCaps, err := data.GetCapabilities("main")
			if err != nil {
				ctx.JSON(Map{"ok": false, "error": err.Error()})
				return
			}
			analyticsCaps, err := data.GetCapabilities("analytics")
			if err != nil {
				ctx.JSON(Map{"ok": false, "error": err.Error()})
				return
			}
			ctx.JSON(Map{"ok": true, "main": mainCaps, "analytics": analyticsCaps})
		},
	})

	infra.Register("data.analytics", http.Router{
		Uri:  "/data/analytics",
		Name: "分析结果",
		Desc: "读取analytics连接的聚合结果",
		Action: func(ctx *http.Context) {
			db := data.Base("analytics")
			defer db.Close()

			items := db.Table("order_agg").Query(Map{
				"$sort": Map{"user_id": ASC},
			})
			if db.Error() != nil {
				ctx.JSON(Map{"ok": false, "error": db.Error().Error()})
				return
			}

			ctx.JSON(Map{"ok": true, "items": items})
		},
	})
}
