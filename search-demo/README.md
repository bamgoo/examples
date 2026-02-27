# search-demo

完整搜索模块示例，包含 `search` 核心 + `file / meilisearch / opensearch / elasticsearch` 驱动。

## 启动

在本目录运行：

```bash
cp config.file.toml config.toml
go run .
```

切换驱动时，把对应配置复制为 `config.toml`：

- `config.meilisearch.toml`
- `config.opensearch.toml`
- `config.elasticsearch.toml`

默认 HTTP 端口：`8100`。

## 路由

```bash
# 查看可用路由
curl http://127.0.0.1:8100/

# 关键词搜索 + 筛选 + 分面
curl "http://127.0.0.1:8100/search?q=search&category=tech&score_gt=8.5&offset=0&limit=10"

# 计数
curl "http://127.0.0.1:8100/search/count?q=search"

# 追加索引
curl -X POST http://127.0.0.1:8100/search/reindex

# 清空索引数据
curl -X POST http://127.0.0.1:8100/search/clear

# 删除文档
curl -X DELETE http://127.0.0.1:8100/search/doc/1005
```

## 说明

1. 启动时会自动创建索引并写入初始化文档。
2. `search.Search(...)` 示例中演示了统一 DSL：`$filters/$sort/$fields/$facets/$highlight`。
3. 示例使用新 DSL：
   - `$filters`: `Map{"category": Map{"$eq": "tech"}, "score": Map{"$gt": 8.5}}`
   - `$sort`: `[]Map{{"score": DESC}, {"id": ASC}}`
4. 不再使用旧格式（`field/desc`、`op/value`）。
5. `search.Upsert(...)` 直接使用 `Map...`，不再使用 `search.Document`。
6. 切换驱动只需要改 TOML，不需要改业务代码。
