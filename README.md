# examples
infrago examples.

## Demos

- `data-demo`: data 模块多驱动示例（sqlite/mysql/postgresql）
- `trace-demo`: trace 模块示例（默认驱动控制台输出）
- `log-demo`: log 模块示例（异步批量写入 + overflow 策略）
- `search-demo`: search 模块 + file/meilisearch/opensearch/elasticsearch 驱动示例（`Upsert(index, rows ...Map)` + `$` DSL）

## Quick Start

```bash
cd search-demo
cp config.file.toml config.toml
go run .
```

```bash
curl "http://127.0.0.1:8100/search?q=search&category=tech&score_gt=8.5&offset=0&limit=10"
```
