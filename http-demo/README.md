# web-demo

单端口多站点示例（`www`、`file`、`api` 三个站点）。

## 运行

默认监听 `0.0.0.0:8090`。

## 测试

用 `--resolve` 模拟不同 Host（不改 hosts 文件）：

```bash
# www site
curl --resolve www.demo.local:8090:127.0.0.1 http://www.demo.local:8090/
curl --resolve www.demo.local:8090:127.0.0.1 http://www.demo.local:8090/about

# file site
curl --resolve file.demo.local:8090:127.0.0.1 http://file.demo.local:8090/list
curl --resolve file.demo.local:8090:127.0.0.1 http://file.demo.local:8090/download/logo.png

# api site
curl --resolve api.demo.local:8090:127.0.0.1 http://api.demo.local:8090/ping
curl --resolve api.demo.local:8090:127.0.0.1 http://api.demo.local:8090/user/1001
curl --resolve api.demo.local:8090:127.0.0.1 http://api.demo.local:8090/fail

# 未配置 site（sys）会落到空站点（.index）
curl --resolve sys.demo.local:8090:127.0.0.1 http://sys.demo.local:8090/

# 直接访问端口也会走空站点
curl http://127.0.0.1:8090/

# 自定义 Found/Error（全站点）
curl --resolve www.demo.local:8090:127.0.0.1 http://www.demo.local:8090/not-found
```

- `*.access` 是全站点 Filter，请求和响应都会打印日志。
- `*.custom` 是全站点 Handler，覆盖 404 和 500 输出格式。
- `.index` 是空站点路由，只在未命中已配置站点（或直接访问 IP/端口）时生效。
