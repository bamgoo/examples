package main

import (
	"fmt"
	"time"

	"github.com/bamgoo/bamgoo"
	. "github.com/bamgoo/base"
	_ "github.com/bamgoo/builtin"
	"github.com/bamgoo/data"
	_ "github.com/bamgoo/data-mysql"
	_ "github.com/bamgoo/data-pgsql"
	_ "github.com/bamgoo/data-sqlite"
	"github.com/bamgoo/http"
)

func main() {
	bamgoo.Go()
}

func init() {
	bamgoo.Register("user", data.Table{
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

	bamgoo.Register("order", data.Table{
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
	bamgoo.Register("order_agg", data.Table{
		Name: "订单聚合快照", Desc: "asdf",
		Schema: "public", Key: "user_id",
		Fields: Vars{
			"user_id": Var{Type: "int"},
			"total":   Var{Type: "float"},
			"cnt":     Var{Type: "int"},
			"updated": Var{Type: "timestamp"},
		},
	})

	bamgoo.Register(bamgoo.START, bamgoo.Trigger{
		Name: "Data Demo", Desc: "Run data module demo on startup",
		Action: func(ctx *bamgoo.Context) {
			db := data.Base("main")
			defer db.Close()
			fmt.Println("capabilities", db.Capabilities())

			adb := data.Base("analytics")
			defer adb.Close()

			_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS "user" (
				"id" INTEGER PRIMARY KEY,
				"name" TEXT NOT NULL,
				"status" TEXT,
				"login_times" INTEGER DEFAULT 0,
				"profile" TEXT,
				"tags" TEXT,
				"created" DATETIME
			)`)

			_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS "order" (
				"id" INTEGER PRIMARY KEY,
				"user_id" INTEGER NOT NULL,
				"amount" REAL,
				"status" TEXT,
				"created" DATETIME
			)`)
			_, _ = adb.Exec(`CREATE TABLE IF NOT EXISTS "order_agg" (
				"user_id" INTEGER PRIMARY KEY,
				"total" REAL,
				"cnt" INTEGER,
				"updated" DATETIME
			)`)

			user, err := db.Table("user").Upsert(Map{
				"$set": Map{
					"name":    "Alice",
					"status":  "active",
					"profile": Map{"city": "Shanghai", "vip": true},
					"tags":    []string{"go", "db", "cloud"},
					"created": time.Now(),
				},
				"$inc": Map{
					"login_times": 1,
				},
			}, Map{"id": 1001})
			if err != nil {
				fmt.Println("upsert user error:", err)
				return
			}

			_, _ = db.Table("order").Create(Map{"id": 1, "user_id": 1001, "amount": 39.5, "status": "paid", "created": time.Now()})
			_, _ = db.Table("order").Create(Map{"id": 2, "user_id": 1001, "amount": 72.0, "status": "paid", "created": time.Now()})
			_, _ = db.Table("order").Create(Map{"id": 3, "user_id": 1001, "amount": 9.0, "status": "new", "created": time.Now()})

			rows, err := db.Table("order").Aggregate(Map{
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
			if err != nil {
				fmt.Println("aggregate error:", err)
				return
			}
			for _, row := range rows {
				_, _ = adb.Table("order_agg").Upsert(Map{
					"$set": Map{
						"total":   row["total"],
						"cnt":     row["cnt"],
						"updated": time.Now(),
					},
				}, Map{"user_id": row["user_id"]})
			}

			tagged, err := db.Table("user").Query(Map{
				"tags":  Map{"$contains": []string{"go"}},
				"$sort": Map{"id": 1},
			})
			if err != nil {
				fmt.Println("contains query error:", err)
			}

			_ = db.Table("order").LimitRange(2, func(item Map) Res {
				fmt.Println("range item", item)
				return bamgoo.OK
			}, Map{
				"$sort": Map{"id": 1},
			})

			fmt.Println("user", user)
			fmt.Println("aggregate rows", rows)
			fmt.Println("tagged users", tagged)
		},
	})

	bamgoo.Register("data.user", http.Router{
		Uri:  "/data/user/{id}",
		Name: "查询用户",
		Desc: "按ID查询用户",
		Action: func(ctx *http.Context) {
			db := data.Base("main")
			defer db.Close()

			user, err := db.Table("user").Entity(ctx.Params["id"])
			if err != nil {
				ctx.JSON(Map{"ok": false, "error": err.Error()})
				return
			}

			ctx.JSON(Map{"ok": true, "user": user})
		},
	})

	bamgoo.Register("data.orders", http.Router{
		Uri:  "/data/orders",
		Name: "订单聚合",
		Desc: "订单聚合统计",
		Action: func(ctx *http.Context) {
			db := data.Base("main")
			defer db.Close()

			rows, err := db.Table("order").Aggregate(Map{
				"status": "paid",
				"$group": []string{"user_id"},
				"$agg": Map{
					"total": Map{"sum": "amount"},
					"cnt":   Map{"count": "*"},
				},
			})
			if err != nil {
				ctx.JSON(Map{"ok": false, "error": err.Error()})
				return
			}

			ctx.JSON(Map{"ok": true, "items": rows})
		},
	})

	bamgoo.Register("data.capabilities", http.Router{
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

	bamgoo.Register("data.analytics", http.Router{
		Uri:  "/data/analytics",
		Name: "分析结果",
		Desc: "读取analytics连接的聚合结果",
		Action: func(ctx *http.Context) {
			db := data.Base("analytics")
			defer db.Close()

			items, err := db.Table("order_agg").Query(Map{
				"$sort": Map{"user_id": 1},
			})
			if err != nil {
				ctx.JSON(Map{"ok": false, "error": err.Error()})
				return
			}

			ctx.JSON(Map{"ok": true, "items": items})
		},
	})
}
