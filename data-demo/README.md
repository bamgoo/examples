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
