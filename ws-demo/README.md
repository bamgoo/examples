# ws-demo

最小可跑的 `web + ws` 示例。

## 功能

- `web.Context.Upgrade()`
- `ws.Hook`
- `ws.Filter`
- `ws.Message`
- `ws.Command`
- `ctx.Reply()`
- `ctx.Broadcast()`
- `ctx.Groupcast()`
- `ctx.BindUser()`
- `ctx.PushUserResult()`
- 本地投递统计 `BroadcastResult / GroupcastResult`
- 协议导出 `/ws/export`
- 运行指标 `/ws/metrics`

## 运行

```bash
cd examples/ws-demo
go run .
```

默认监听 `http://127.0.0.1:8080/`。

页面里可以直接测试：

- echo
- join / groupcast
- broadcast
- bind user / push user
