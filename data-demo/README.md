# data-demo

Run in this folder with sqlite:

```bash
cp config.sqlite.toml config.toml
go run .
```

Switch to pgsql/mysql by copying corresponding config file to `config.toml`.

HTTP demo routes (default port `8099`):

```bash
curl http://127.0.0.1:8099/data/user/1001
curl http://127.0.0.1:8099/data/orders
curl http://127.0.0.1:8099/data/capabilities
curl http://127.0.0.1:8099/data/analytics
```

Write API semantics used in demo:

```go
// single row (default sort by primary key asc when no $sort)
updated := db.Table("order").Update(base.Map{
  "$set": base.Map{"status": "processing"},
}, base.Map{"status": "new"})

deleted := db.Table("order").Delete(base.Map{"status": "new"})

// all matched rows
affected := db.Table("order").UpdateMany(base.Map{
  "$set": base.Map{"status": "paid"},
}, base.Map{"status": "processing"})

affected = db.Table("order").DeleteMany(base.Map{
  "status": "paid",
  "id": base.Map{"$gt": 2},
})
_, _, _ = updated, deleted, affected
```
