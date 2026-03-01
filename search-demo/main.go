package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/infrago/infra"
	. "github.com/infrago/base"
	_ "github.com/infrago/builtin"
	"github.com/infrago/http"
	"github.com/infrago/search"

	_ "github.com/infrago/search-elasticsearch"
	_ "github.com/infrago/search-file"
	_ "github.com/infrago/search-meilisearch"
	_ "github.com/infrago/search-opensearch"
)

const indexName = "article"

func main() {
	infra.Go()
}

func init() {
	infra.Register(indexName, search.Index{
		Name:        "文章索引",
		Primary:     "id",
		StrictWrite: true,
		StrictRead:  false,
		Attributes: Vars{
			"id":       Var{Type: "string", Required: true},
			"title":    Var{Type: "string", Required: true},
			"content":  Var{Type: "string"},
			"category": Var{Type: "string"},
			"tags":     Var{Type: "[string]"},
			"score":    Var{Type: "float"},
			"created":  Var{Type: "timestamp"},
		},
	})

	infra.Register(infra.START, infra.Trigger{
		Name: "search-demo-init",
		Desc: "create index and seed docs",
		Action: func(ctx *infra.Context) {
			_ = search.Upsert(indexName,
				Map{"id": "1001", "title": "Go 微服务实战", "content": "Bamgoo + Nomad + NATS 快速搭建高性能服务", "category": "tech", "tags": []string{"go", "microservice", "nomad"}, "score": 9.7, "created": time.Now().Unix()},
				Map{"id": "1002", "title": "搜索系统设计", "content": "全文检索、过滤、分面和高亮设计要点", "category": "arch", "tags": []string{"search", "design"}, "score": 9.3, "created": time.Now().Unix()},
				Map{"id": "1003", "title": "Meilisearch 上手", "content": "轻量搜索服务快速接入指南", "category": "tech", "tags": []string{"meilisearch", "search"}, "score": 8.8, "created": time.Now().Unix()},
				Map{"id": "1004", "title": "OpenSearch 聚合", "content": "通过 terms 聚合做 category 统计", "category": "arch", "tags": []string{"opensearch", "aggs"}, "score": 8.9, "created": time.Now().Unix()},
			)
		},
	})

	infra.Register("search.index", http.Router{
		Uri:  "/",
		Name: "search-demo-index",
		Desc: "help",
		Action: func(ctx *http.Context) {
			ctx.JSON(Map{
				"ok": true,
				"routes": []string{
					"GET /search?q=go&category=tech&offset=0&limit=10",
					"GET /search/count?q=search",
					"POST /search/reindex",
					"POST /search/clear",
					"DELETE /search/doc/{id}",
				},
			})
		},
	})

	infra.Register("search.query", http.Router{
		Uri:  "/search",
		Name: "search-query",
		Desc: "query docs",
		Action: func(ctx *http.Context) {
			q, _ := ctx.Query["q"].(string)
			category, _ := ctx.Query["category"].(string)
			scoreGt := toFloat(ctx.Query["score_gt"], 0)
			offset := toInt(ctx.Query["offset"], 0)
			limit := toInt(ctx.Query["limit"], 10)

			filters := Map{}
			if strings.TrimSpace(category) != "" {
				filters["category"] = Map{"$eq": category}
			}
			if scoreGt > 0 {
				filters["score"] = Map{"$gt": scoreGt}
			}

			res, err := search.Search(indexName, q, Map{
				"$offset":    offset,
				"$limit":     limit,
				"$filters":   filters,
				"$sort":      []Map{{"score": DESC}, {"id": ASC}},
				"$fields":    []string{"title", "content", "category", "tags", "score", "created"},
				"$facets":    []string{"category"},
				"$highlight": []string{"title", "content"},
			})
			if err != nil {
				ctx.JSON(Map{"ok": false, "error": err.Error()})
				return
			}
			ctx.JSON(Map{"ok": true, "query": q, "result": res})
		},
	})

	infra.Register("search.count", http.Router{
		Uri:  "/search/count",
		Name: "search-count",
		Desc: "count docs",
		Action: func(ctx *http.Context) {
			q, _ := ctx.Query["q"].(string)
			total, err := search.Count(indexName, q)
			if err != nil {
				ctx.JSON(Map{"ok": false, "error": err.Error()})
				return
			}
			ctx.JSON(Map{"ok": true, "query": q, "total": total})
		},
	})

	infra.Register("search.reindex", http.Router{
		Uri:  "/search/reindex",
		Name: "search-reindex",
		Desc: "reindex docs",
		Action: func(ctx *http.Context) {
			err := search.Upsert(indexName,
				Map{"id": "1005", "title": "Elasticsearch 实战", "content": "字段过滤与分面查询示例", "category": "tech", "tags": []string{"es", "filter", "facet"}, "score": 9.1, "created": time.Now().Unix()},
			)
			if err != nil {
				ctx.JSON(Map{"ok": false, "error": err.Error()})
				return
			}
			ctx.JSON(Map{"ok": true, "msg": "reindex done"})
		},
	})

	infra.Register("search.delete", http.Router{
		Uri:  "/search/doc/{id}",
		Name: "search-delete",
		Desc: "delete doc",
		Action: func(ctx *http.Context) {
			id, _ := ctx.Params["id"].(string)
			if strings.TrimSpace(id) == "" {
				ctx.JSON(Map{"ok": false, "error": "id is required"})
				return
			}
			if err := search.Delete(indexName, []string{id}); err != nil {
				ctx.JSON(Map{"ok": false, "error": err.Error()})
				return
			}
			ctx.JSON(Map{"ok": true, "id": id})
		},
	})

	infra.Register("search.clear", http.Router{
		Uri:  "/search/clear",
		Name: "search-clear",
		Desc: "clear index",
		Action: func(ctx *http.Context) {
			if err := search.Clear(indexName); err != nil {
				ctx.JSON(Map{"ok": false, "error": err.Error()})
				return
			}
			ctx.JSON(Map{"ok": true, "index": indexName})
		},
	})
}

func toInt(v Any, def int) int {
	switch vv := v.(type) {
	case int:
		return vv
	case int64:
		return int(vv)
	case float64:
		return int(vv)
	case string:
		if vv == "" {
			return def
		}
		out, err := strconv.Atoi(vv)
		if err == nil {
			return out
		}
	}
	return def
}

func toFloat(v Any, def float64) float64 {
	switch vv := v.(type) {
	case float64:
		return vv
	case float32:
		return float64(vv)
	case int:
		return float64(vv)
	case int64:
		return float64(vv)
	case string:
		if vv == "" {
			return def
		}
		out, err := strconv.ParseFloat(vv, 64)
		if err == nil {
			return out
		}
	}
	return def
}
