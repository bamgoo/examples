# web-multisite

单端口多站点（子域名）示例。

## 运行

```bash
cd /Users/yatniel/coding/bamgoo/examples/web-multisite
go run .
```

默认监听 `0.0.0.0:8090`。

## 测试

不改系统 hosts 的方式（推荐临时验证）：

```bash
curl --resolve www.bamgoo.local:8090:127.0.0.1 http://www.bamgoo.local:8090/
curl --resolve www.bamgoo.local:8090:127.0.0.1 http://www.bamgoo.local:8090/about

curl --resolve user.bamgoo.local:8090:127.0.0.1 http://user.bamgoo.local:8090/
curl --resolve user.bamgoo.local:8090:127.0.0.1 http://user.bamgoo.local:8090/profile/1001

curl --resolve sys.bamgoo.local:8090:127.0.0.1 http://sys.bamgoo.local:8090/health
```

`*.logger` 会自动挂到 `www/user/sys` 三个站点。
